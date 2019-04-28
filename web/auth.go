package web

import (
	"drive/drive/storage"
	"drive/errors"
	"drive/session"
	"drive/templates"
	"fmt"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		user, err := storage.Authenticate(r.PostFormValue("username"), r.PostFormValue("password"))
		if err != nil {
			errors.Error(w, r, errors.Augment(err, errors.BadCredentials, "Could not validate given credentials"))
			return
		}
		fmt.Println("=>", user, r.PostFormValue("username"), r.PostFormValue("password"))

		session.SetSessionUser(r, w, user)
		if err != nil {
			errors.Error(w, r, errors.Augment(err, errors.Session, "Could not store User in session"))
			return
		}

		if next, ok := r.URL.Query()["next"]; ok && len(next) >= 1 {
			http.Redirect(w, r, next[0], 302)
		}
		http.Redirect(w, r, "/", 302)
	}

	err := templates.Render(w, http.StatusOK, "login", nil)
	if err != nil {
		fmt.Println("login render error: ", err)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	sess, _ := session.GetSession(r, w)
	sess.Clear()
	sess.Save(r, w)
	http.Redirect(w, r, "/login", 302)
}
