package main

import (
	"drive/fs"
	"drive/session"
	"fmt"
	"html/template"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {

	sess := session.GetSession(r, w)

	if r.Method == http.MethodPost {

		username, password := r.PostFormValue("username"), r.PostFormValue("password")
		if username == "" || password == "" {
			http.Error(w, "empty", http.StatusBadRequest)
			return
		}
		if username == "matt" && password == "mehmet" {
			sess.Set("authenticated", true)
			sess.Set("username", "matt")
		} else {
			sess.Set("authenticated", false)
		}
		sess.Save()
		http.Redirect(w, r, "/secret", 302)

	}
	t, _ := template.ParseFiles("templates/login.html")
	t.Execute(w, nil)

}

func logout(w http.ResponseWriter, r *http.Request) {
	sess := session.GetSession(r, w)

	// Revoke users authentication
	sess.Set("authenticated", false)
	sess.Save()
}

func Index(w http.ResponseWriter, r *http.Request) {
	auth := make(map[int]int) //session.Get(r, "authenticated")
	fmt.Println(auth)
	if auth == nil {
		http.Redirect(w, r, "/login", 301)

		return
	}
	fmt.Println("Index")

	//	username := session.Get(r, "username").(string)
	//	storage := fs.GetStorage(username, path)
	//	file := storage.GetFile(path)

	//	views.Render(file, username)

	var storage = &fs.FileSystemStorage{Location: "/Users/mi/go"}
	storage.ServeHTTP(w, r)
}

func secret(w http.ResponseWriter, r *http.Request) {
	auth := session.Get(r, "authenticated")
	username := session.Get(r, "username").(string)
	if auth == nil || !auth.(bool) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Print secret message
	fmt.Fprintln(w, username)
}
