package session

import (
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var SESSION_KEY = "aslkdjfheqhlieufhkasjbd"
var SESSION_NAME = "gosessionid"

// store will hold all session data
var store *sessions.CookieStore

func init() {

	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 24,
		HttpOnly: true,
	}

}

type Session struct {
	Session *sessions.Session
	r       *http.Request
	w       http.ResponseWriter
}

// Clear the current session
func (s *Session) Clear() {
	for k := range s.Session.Values {
		s.Delete(k)
	}
}

// Delete a value from the current session.
func (s *Session) Delete(name interface{}) {
	delete(s.Session.Values, name)
}

func (s *Session) Get(name interface{}) interface{} {
	return s.Session.Values[name]
}

// GetOnce gets a value from the current session and then deletes it.
func (s *Session) GetOnce(name interface{}) interface{} {
	if x, ok := s.Session.Values[name]; ok {
		s.Delete(name)
		return x
	}
	return nil
}

func (s *Session) Save() error {
	return s.Session.Save(s.r, s.w)
}

// Set a value onto the current session. If a value with that name
// already exists it will be overridden with the new value.
func (s *Session) Set(name, value interface{}) {
	s.Session.Values[name] = value
}

// Get a session using a request and response.
func GetSession(r *http.Request, w http.ResponseWriter) *Session {

	session, err := store.Get(r, SESSION_NAME)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	return &Session{
		Session: session,
		r:       r,
		w:       w,
	}
}

func Get(r *http.Request, name interface{}) interface{} {
	session := GetSession(r, nil)
	if session == nil {
		return nil
	}

	return session.Get(name)
}

func Set(r *http.Request, key, value interface{}) error {
	session := GetSession(r, nil)
	session.Set(key, value)
	return nil
}
