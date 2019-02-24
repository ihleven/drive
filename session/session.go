package session

import (
	"drive/auth"
	"encoding/gob"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

var SESSION_KEY = "aslkdjfheqhlieufhkasjbd"
var SESSION_NAME = "gosessionid"

// store will hold all session data
var store = sessions.NewCookieStore([]byte("something-very-secret"))

func init() {

	//authKeyOne := securecookie.GenerateRandomKey(64)
	//encryptionKeyOne := securecookie.GenerateRandomKey(32)
	//store = sessions.NewCookieStore(
	//	[]byte("asdaskdhasdhgsajdgasdsadksakdhasidoajsdousahdopj"),
	//	[]byte("hhfjhtdzjtfkhgkjfkufkztfjztfkuztfkztdhtesrgesdjg"),
	//)
	store.Options = &sessions.Options{
		//Domain:   "localhost",
		Path:     "/",
		MaxAge:   3600 * 8, // 8 hours
		HttpOnly: true,
		//Secure:   false,
	}
	gob.Register(auth.Account{})
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
func (s *Session) GetString(name string) string {
	str := s.Session.Values[name].(string)
	fmt.Println("---", str)
	return str
}

// GetOnce gets a value from the current session and then deletes it.
func (s *Session) GetOnce(name interface{}) interface{} {
	if x, ok := s.Session.Values[name]; ok {
		s.Delete(name)
		return x
	}
	return nil
}

func (s *Session) Save(r *http.Request, w http.ResponseWriter) error {
	return s.Session.Save(r, w)
}

// Set a value onto the current session. If a value with that name
// already exists it will be overridden with the new value.
func (s *Session) Set(name, value interface{}) {
	s.Session.Values[name] = value
}

// Get a session using a request and response.
func GetSession(r *http.Request, w http.ResponseWriter) (*Session, error) {

	session, err := store.Get(r, SESSION_NAME)
	if err != nil {
		fmt.Println("ERROR session.GetSession", err.Error())
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	return &Session{
		Session: session,
		r:       r,
		w:       w,
	}, nil
}

func Get(r *http.Request, name interface{}) interface{} {
	session, _ := GetSession(r, nil)
	if session == nil {
		return nil
	}
	var value interface{}
	value = session.Get(name)
	fmt.Println("VALUE:", value)
	return value
}
func GetString(r *http.Request, name interface{}) string {
	session, _ := GetSession(r, nil)
	if session == nil {
		return ""
	}
	var value interface{}
	value = session.Get(name)
	if value != nil {
		return value.(string)
	}
	return ""
}

func Set(r *http.Request, key, value interface{}) error {
	session, _ := GetSession(r, nil)
	session.Set(key, value)
	return nil
}

func GetSessionUser(r *http.Request, w http.ResponseWriter) (*auth.Account, error) {
	sess, err := GetSession(r, w)
	if err != nil {
		return nil, err
	}

	val := sess.Get("user")
	var user = auth.Account{}
	user, ok := val.(auth.Account)
	if !ok {
		return &auth.Account{Authenticated: false}, nil
	}
	return &user, nil
}

func SetSessionUser(r *http.Request, w http.ResponseWriter, user *auth.Account) (err error) {
	sess, err := GetSession(r, w)
	if err != nil {
		fmt.Println("SetSessoinUser GetSession Error: ", err)
		return
	}
	sess.Set("user", user)
	sess.Save(r, w)
	//	err = session.Save()
	return
}

func AuthUser(r *http.Request, w http.ResponseWriter) (*auth.Account, error) {
	sess, err := GetSession(r, w)
	if err != nil {
		return nil, err
	}

	val := sess.Get("user")
	var user = auth.Account{}
	user, ok := val.(auth.Account)
	if !ok {
		return &auth.Account{Authenticated: false}, nil
	}
	return &user, nil
}
