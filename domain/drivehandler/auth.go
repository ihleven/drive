package drivehandler

import (
	"drive/domain/usecase"
	"drive/session"
	"drive/web"
	"fmt"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		user, err := usecase.Authenticate(r.PostFormValue("username"), r.PostFormValue("password"))
		if err != nil {
			web.ErrorResponder(w, err.Error(), 500)
			return
		}
		fmt.Println("=>", user, r.PostFormValue("username"), r.PostFormValue("password"))

		session.SetSessionUser(r, w, user)
		if err != nil {
			web.ErrorResponder(w, "session error"+err.Error(), 500)
			return
		}

		if next, ok := r.URL.Query()["next"]; ok && len(next) >= 1 {
			http.Redirect(w, r, next[0], 302)
		}
		http.Redirect(w, r, "/", 302)
	}

	err := rnd.HTML(w, http.StatusOK, "login", nil)
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
