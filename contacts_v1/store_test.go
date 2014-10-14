package main

import (
	"io/ioutil"
	"log"
	"reflect"
	"testing"
	"time"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func isRedisRunning() bool {
	return IsProcessRunning("redis-server")
}

func TestRoundTripInMemorySessionStore(t *testing.T) {
	age, purgeInterval := uint(1), uint(1)
	store := NewInMemorySessionStore(age, purgeInterval)

	RunRoundtripSessionStoreTest(t, store)
}

func TestInMemorySessionStoreRecordNotFound(t *testing.T) {
	age, purgeInterval := uint(1), uint(1)
	store := NewInMemorySessionStore(age, purgeInterval)

	RunSessionStoreRecordNotFoundTest(t, store)
}

func TestRoundTripRedisSessionStore(t *testing.T) {
	if !isRedisRunning() {
		t.Skip("No redis instance running")
	}

	age := uint(15)

	pool := NewRedisPool(":6379", "")
	defer pool.Close()

	store := NewRedisSessionStore(pool, age)

	RunRoundtripSessionStoreTest(t, store)
}

func TestRedisSessionStoreRecordNotFound(t *testing.T) {
	if !isRedisRunning() {
		t.Skip("No redis instance running")
	}

	age := uint(15)

	pool := NewRedisPool(":6379", "")
	defer pool.Close()

	store := NewRedisSessionStore(pool, age)

	RunSessionStoreRecordNotFoundTest(t, store)
}

func TestRoundTripInMemoryUserStore(t *testing.T) {
	store := NewInMemoryUserStore()

	RunRoundtripUserStoreTest(t, store)
}

func TestInMemoryUserStoreRecordNotFound(t *testing.T) {
	store := NewInMemoryUserStore()

	RunUserStoreRecordNotFoundTest(t, store)
}

func TestRoundtripRedisUserStore(t *testing.T) {
	if !isRedisRunning() {
		t.Skip("No redis instance running")
	}

	pool := NewRedisPool(":6379", "")
	defer pool.Close()

	store := NewRedisUserStore(pool)

	RunRoundtripUserStoreTest(t, store)
}

func TestRedisUserStoreRecordNotFound(t *testing.T) {
	if !isRedisRunning() {
		t.Skip("No redis instance running")
	}

	pool := NewRedisPool(":6379", "")
	defer pool.Close()

	store := NewRedisUserStore(pool)
	RunUserStoreRecordNotFoundTest(t, store)
}

/*
Helper functions
*/
func RunRoundtripSessionStoreTest(t *testing.T, store SessionStore) {
	spec := &Spec{t}

	original := &Session{
		Id:         "s100",
		UserName:   "Ted",
		Data:       map[string]interface{}{"A1": 1, "A2": "...."},
		LastAccess: time.Now(),
	}

	err := store.Save(original)
	spec.Assert(err == nil, "Unexpected error : %s", err)

	retrieved, err := store.Get(original.Id)
	spec.Assert(err == nil, "Unexpected error : %s", err)

	spec.Assert(reflect.DeepEqual(original, retrieved), "Expected [%v] but got [%v]", original, retrieved)
}

func RunSessionStoreRecordNotFoundTest(t *testing.T, store SessionStore) {
	spec := &Spec{t}

	retrieved, err := store.Get("DOESNOTEXIST")

	spec.Assert(err != nil, "Expected error")
	spec.Assert(retrieved == nil, "Expected session to be nil")
}

func RunRoundtripUserStoreTest(t *testing.T, store UserStore) {
	spec := &Spec{t}

	original := &User{
		Id:        "pmcgrath",
		FirstName: "Pat",
		LastName:  "Mc Grath",
		Email:     "pmcgrath@gmail.com",
		Password:  "pass",
		Contacts: []Contact{
			Contact{
				Id:        "pmcgrath",
				FirstName: "Peter",
				LastName:  "Mc Grath",
				Phones: []Phone{
					Phone{
						Description: "Home",
						Number:      "44 066 7132310",
					},
				},
			},
		},
	}

	err := store.Save(original)
	spec.Assert(err == nil, "Unexpected error : %s", err)

	retrieved, err := store.Get(original.Id)
	spec.Assert(err == nil, "Unexpected error : %s", err)

	spec.Assert(reflect.DeepEqual(original, retrieved), "Expected [%v] but got [%v]", original, retrieved)
}

func RunUserStoreRecordNotFoundTest(t *testing.T, store UserStore) {
	spec := &Spec{t}

	retrieved, err := store.Get("DoesNotExist")

	spec.Assert(err != nil, "Expected error")
	spec.Assert(retrieved == nil, "Expected user to be nil")
}
