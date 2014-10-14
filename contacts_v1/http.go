package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

/*
Our http handler func function signature
*/
type ContextualHandlerFunc func(http.ResponseWriter, *http.Request, *RequestContext)

/*
Our http handler interface
*/
type ContextualHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request, *RequestContext)
}

/*
Request context type
*/
type RequestContext struct {
	Id        string
	StartTime time.Time
	Session   *Session
	Data      map[string]interface{}
}

func (c *RequestContext) GetSessionId() string {
	if c.Session != nil {
		return c.Session.Id
	}

	return ""
}

func (c *RequestContext) GetUserName() string {
	if c.Session != nil {
		return c.Session.UserName
	}

	return ""
}

func (c *RequestContext) IsLoggedIn() bool {
	return c.GetUserName() != ""
}

func (c *RequestContext) GetLogMessagePrefix() string {
	return fmt.Sprintf("%s %s [%s]", c.Id, c.GetSessionId(), c.GetUserName())
}

/*
Middleware types
*/
/*
Init function which returns a http.HandlerFunc wrapper which we can use when calling http.HandleFunc in main
*/
func CreateInitHandlerFunc(next ContextualHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := &RequestContext{
			Id:        Uuid(),
			StartTime: time.Now(),
			Data:      make(map[string]interface{}, 0),
		}

		next.ServeHTTP(w, r, c)
	}
}

/*
Logging middleware
*/
type SpyResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (w *SpyResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	w.StatusCode = code
}

func NewSpyResponseWriter(inner http.ResponseWriter) *SpyResponseWriter {
	return &SpyResponseWriter{
		ResponseWriter: inner,
		StatusCode:     200,
	}
}

type LoggingHandler struct {
	Next ContextualHandler
}

func (h *LoggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, c *RequestContext) {
	spyResponseWriter := NewSpyResponseWriter(w)

	h.Next.ServeHTTP(spyResponseWriter, r, c)

	log.Printf("%s %s %s %d\n", c.GetLogMessagePrefix(), r.URL.Path, r.Method, spyResponseWriter.StatusCode)
}

func NewLoggingHandler(next ContextualHandler) ContextualHandler {
	return &LoggingHandler{Next: next}
}

/*
Session middleware
*/
type SessionHandler struct {
	Store SessionStore
	Next  ContextualHandler
}

func (h *SessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, c *RequestContext) {
	var s *Session
	cookie, err := r.Cookie("SessionId")
	if err != nil {
		s = &Session{
			Id: Uuid(),
		}
	} else {
		s, _ = h.Store.Get(cookie.Value)
		if s == nil {
			s = &Session{
				Id: Uuid(),
			}
		}
	}
	c.Session = s

	// *** Need to write cookie before we pass to the next handler as we will not succed in writing the cookie after the next handler has completed, if it has written some contentto the repsonse body
	// Rewrite session cookie
	cookie = &http.Cookie{
		Name:     "SessionId",
		Value:    c.Session.Id,
		Path:     "/",
		Domain:   "", // Chrome will not include if value is "localhost" the Cookie header in requests, seems to need 2 dots see http://stackoverflow.com/questions/21865681/sessions-variables-in-golang-not-saved-while-using-gorilla-sessions
		MaxAge:   int(h.Store.GetAge()),
		Secure:   false, // Requires TLS to make true
		HttpOnly: true,
	}
	// Only get to write one cookie, so this will overwrite any existing cookies
	http.SetCookie(w, cookie)

	h.Next.ServeHTTP(w, r, c)

	h.Store.Save(s)
}

func NewSessionHandler(store SessionStore, next ContextualHandler) ContextualHandler {
	return &SessionHandler{Store: store, Next: next}
}

/*
Authorisation middleware
*/
type AuthorisationHandler struct {
	Next ContextualHandler
}

func (h *AuthorisationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, c *RequestContext) {
	if !c.IsLoggedIn() {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	h.Next.ServeHTTP(w, r, c)
}

func NewAuthorisationHandler(next ContextualHandler) ContextualHandler {
	return &AuthorisationHandler{Next: next}
}
