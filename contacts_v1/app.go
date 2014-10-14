package main

import (
	_ "expvar" // So we can access debug/vars
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func init() {
	pid := os.Getpid()
	log.SetPrefix(fmt.Sprintf("%d ", pid))
}

func openStores() (sessionStore SessionStore, userStore UserStore) {
	redisAddress := GetOrDefaultEnv("REDIS_ADDRESS", "")
	redisPassword := GetOrDefaultEnv("REDIS_PASSWORD", "")
	sessionTimeoutInMinutes, _ := strconv.Atoi(GetOrDefaultEnv("WEBAPP_SESSION_TIMEOUT_IN_MINUTES", "20"))

	sessionTimeoutInSeconds := uint(sessionTimeoutInMinutes * 60)

	if redisAddress != "" {
		log.Printf("Using redis stores %s\n", redisAddress)
		pool := NewRedisPool(redisAddress, redisPassword)

		sessionStore = NewRedisSessionStore(pool, sessionTimeoutInSeconds)
		userStore = NewRedisUserStore(pool)
	} else {
		log.Println("Using in memory stores - will add 'pmcgrath' user")

		sessionStore = NewInMemorySessionStore(sessionTimeoutInSeconds, sessionTimeoutInSeconds) // Purge and timeout are same value
		userStore = NewInMemoryUserStore()

		// Add a user so we have a user to work with
		userStore.Save(&User{
			Id:        "pmcgrath",
			FirstName: "Pat",
			LastName:  "Mc Grath",
			Email:     "pmcgrat@gmail.com",
			Password:  "pass",
			Contacts:  make([]Contact, 0),
		})
	}

	return
}

func closeStores(sessionStore SessionStore, userStore UserStore) {
	// If redis stores, close redis pool - same pool shared by both stores
	if store, ok := sessionStore.(*RedisSessionStore); ok {
		log.Println("Closing redis pool")
		store.pool.Close()
	}
}

func main() {
	webAppAddress := GetOrDefaultEnv("WEBAPP_ADDRESS", ":8080")

	sessionStore, userStore := openStores()
	defer closeStores(sessionStore, userStore)

	rootHandler := &RootHandler{}
	assetsHandler := &AssetsHandler{}
	contactApiHandler := &ContactApiHandler{PathPrefix: "/api/v1/contacts/", Store: userStore}
	contactsApiHandler := &ContactsApiHandler{PathPrefix: "/api/v1/contacts/", Store: userStore}
	logInApiHandler := &LogInApiHandler{Store: userStore}

	router := NewRouter()
	router.Add(`^/?$`, rootHandler)
	router.Add(`^/assets/.*`, assetsHandler)
	router.Add(`^/api/v1/contacts/[\w-]{5,36}/[\w-]{5,36}/?$`, contactApiHandler)
	router.Add(`^/api/v1/contacts/[\w-]{5,36}/?$`, contactsApiHandler)
	router.Add(`^/api/v1/login/?$`, logInApiHandler)

	http.Handle("/", CreateInitHandlerFunc(NewLoggingHandler(NewSessionHandler(sessionStore, router))))                                 // Don't need to be an authenticated user
	http.Handle("/assets/", CreateInitHandlerFunc(NewLoggingHandler(router)))                                                           // Don't need a session
	http.Handle("/api/v1/", CreateInitHandlerFunc(NewLoggingHandler(NewSessionHandler(sessionStore, NewAuthorisationHandler(router))))) // Must be an authenticated user
	http.Handle("/api/v1/login", CreateInitHandlerFunc(NewLoggingHandler(NewSessionHandler(sessionStore, router))))                     // Subset of api that does not need to be an authenticated user, this is a single exception, if we move log in\out out of api we can avoid this

	log.Printf("Started, listening on %s\n", webAppAddress)
	http.ListenAndServe(webAppAddress, nil)
}
