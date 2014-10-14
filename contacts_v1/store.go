package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
)

/*
Store interfaces
*/
type SessionStore interface {
	Get(string) (*Session, error)
	Save(*Session) error
	GetAge() uint
}

type UserStore interface {
	Get(id string) (*User, error)
	Save(user *User) error
}

/*
In memory session store
*/
type InMemorySessionStore struct {
	mutex *sync.RWMutex
	age   uint
	data  map[string]*Session
}

func (store *InMemorySessionStore) Get(id string) (*Session, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	s, ok := store.data[id]
	if !ok {
		return nil, fmt.Errorf("Session for id %s not found", id)
	}

	return s, nil
}

func (store *InMemorySessionStore) Save(s *Session) error {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	s.LastAccess = time.Now()
	store.data[s.Id] = s

	return nil
}

func (store *InMemorySessionStore) Purge() {
	log.Println("Purging session store")
	store.mutex.Lock()
	defer store.mutex.Unlock()

	now := time.Now()
	ageAsDuration := time.Duration(store.age) * time.Second
	for id, session := range store.data {
		if session.LastAccess.Add(ageAsDuration).Before(now) {
			log.Printf("Purging session with Id [%s]\n", id)
			delete(store.data, id)
		}
	}
	log.Printf("Session store purge completed, %d session(s) still in store\n", len(store.data))
}

func (store *InMemorySessionStore) GetAge() uint {
	return store.age
}

func NewInMemorySessionStore(age, purgeInterval uint) *InMemorySessionStore {
	store := &InMemorySessionStore{
		mutex: new(sync.RWMutex),
		age:   age,
		data:  make(map[string]*Session, 0),
	}
	go func() {
		for {
			time.Sleep(time.Duration(purgeInterval) * time.Second)
			store.Purge()
		}
	}()

	return store
}

/*
Redis user store - needs to serialize using gob rather than json due to the sessions data being in a map with a interface{} value
*/
type RedisSessionStore struct {
	pool *redis.Pool
	age  uint
}

func (store *RedisSessionStore) Get(id string) (*Session, error) {
	conn := store.pool.Get()
	defer conn.Close()

	redisKey := "session:" + id
	sessionData, err := redis.Bytes(conn.Do("GET", redisKey))
	if err != nil {
		return nil, err
	}

	sessionDataBuffer := bytes.NewBuffer(sessionData)
	decoder := gob.NewDecoder(sessionDataBuffer)
	session := &Session{}
	if err := decoder.Decode(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (store *RedisSessionStore) Save(session *Session) error {
	conn := store.pool.Get()
	defer conn.Close()

	sessionDataBuffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(sessionDataBuffer)
	if err := encoder.Encode(session); err != nil {
		return err
	}

	redisKey := "session:" + session.Id
	_, err := conn.Do("SETEX", redisKey, store.age, sessionDataBuffer)
	if err != nil {
		return err
	}

	return nil
}

func (store *RedisSessionStore) GetAge() uint {
	return store.age
}

func NewRedisSessionStore(pool *redis.Pool, age uint) *RedisSessionStore {
	return &RedisSessionStore{
		pool: pool,
		age:  age,
	}
}

/*
In memory user store
*/
type InMemoryUserStore struct {
	mutex *sync.RWMutex
	data  map[string]*User
}

func (store *InMemoryUserStore) Get(id string) (*User, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	user, ok := store.data[id]
	if !ok {
		return nil, fmt.Errorf("Record not found for [%s]", id)
	}

	return user, nil
}

func (store *InMemoryUserStore) Save(user *User) error {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	store.data[user.Id] = user
	return nil
}

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		mutex: new(sync.RWMutex),
		data:  make(map[string]*User),
	}
}

/*
Redis user store
*/
type RedisUserStore struct {
	pool *redis.Pool
}

func (store *RedisUserStore) Get(id string) (*User, error) {
	conn := store.pool.Get()
	defer conn.Close()

	redisKey := "user:" + id
	values, err := redis.Values(conn.Do("HGETALL", redisKey))
	if err != nil {
		return nil, err
	}

	var data struct {
		FirstName, LastName, Email, Password, ContactsAsJson string
	}
	if err = redis.ScanStruct(values, &data); err != nil {
		return nil, err
	}

	if (data.FirstName + data.LastName + data.Email + data.Password + data.ContactsAsJson) == "" {
		// No data, so we presume no user
		return nil, fmt.Errorf("Record not found for [%s]", redisKey)
	}

	var contacts []Contact
	if data.ContactsAsJson != "" {
		if err = json.Unmarshal([]byte(data.ContactsAsJson), &contacts); err != nil {
			return nil, err
		}
	}

	user := &User{
		Id:        id,
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Password:  data.Password,
		Contacts:  contacts,
	}

	return user, nil
}

func (store *RedisUserStore) Save(user *User) error {
	conn := store.pool.Get()
	defer conn.Close()

	contactsAsJson, err := json.Marshal(user.Contacts)
	if err != nil {
		return err
	}

	redisKey := "user:" + user.Id
	_, err = conn.Do("HMSET", redisKey,
		"FirstName", user.FirstName,
		"LastName", user.LastName,
		"Email", user.Email,
		"Password", user.Password,
		"ContactsAsJson", contactsAsJson)
	if err != nil {
		return err
	}

	return nil
}

func NewRedisUserStore(pool *redis.Pool) *RedisUserStore {
	return &RedisUserStore{pool: pool}
}

/*
Redis pool creation function
*/
func NewRedisPool(redisAddress, redisPassword string) *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", redisAddress)
			if err != nil {
				return nil, err
			}
			if redisPassword != "" {
				if _, err = conn.Do("AUTH", redisPassword); err != nil {
					conn.Close()
					return nil, err
				}
			}
			return conn, nil
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}
}
