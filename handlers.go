package main

import (
	"drive/auth"
	"drive/controller"
	"drive/file"
	"drive/fs"
	"drive/session"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"strings"
)

// Authorization Key
//var authKey = []byte("somesecret")

// Encryption Key
//var encKey = []byte("someothersecret")

func login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		user, err := auth.Authenticate(r.PostFormValue("username"), r.PostFormValue("password"))
		if err != nil {
			//msg, code := toHTTPError(err)
			//http.Error(w, msg, code)
			return
		}
		session.SetSessionUser(r, w, user)
		fmt.Println(user)

		if err != nil {
			fmt.Println("saving error", err)
		}
		//store.Save(req,res,sessionNew)

		http.Redirect(w, r, "/src/drive", 302)

	}
	t, _ := template.ParseFiles("templates/login.html")
	t.Execute(w, nil)

}

func logout(w http.ResponseWriter, r *http.Request) {
	sess, _ := session.GetSession(r, w)
	sess.Clear()
	sess.Save(r, w)
	http.Redirect(w, r, "/login", 302)
}

var storage = &file.FileSystemStorage{Root: "/Users/mi/go"}

func Serve(w http.ResponseWriter, r *http.Request) {

	authuser, err := session.AuthUser(r, w)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	path := strings.TrimPrefix(path.Clean(r.URL.Path), "/serve/")

	file, err := storage.Open(path, os.O_RDONLY, authuser.Uid, authuser.Gid)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()

	if file.IsDir() {
		r.URL.Path = fmt.Sprintf("/serve/%s/index.html", path)
		Serve(w, r)
		return
	}
	http.ServeContent(w, r, file.Name(), file.ModTime(), file.Descriptor)
}

func PathHandler(w http.ResponseWriter, r *http.Request) {

	usr, err := session.GetSessionUser(r, w)

	storage := &file.FileSystemStorage{Root: "/Users/mi/go"}
	f, err := storage.Open(path.Clean(r.URL.Path), os.O_RDONLY, usr.Uid, usr.Gid)
	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)

		return
	}
	file := &file.File{Info: f, Path: path.Clean(r.URL.Path)}
	defer file.Close()

	mode := f.Mode()
	switch {
	case mode.IsRegular():

		controller := controller.FileController{file, usr}
		controller.Render(w, r)

	case mode.IsDir():

		controller := controller.DirController{file, usr}
		controller.Render(w, r)

	case mode&os.ModeSymlink != 0:
		fmt.Println("symbolic link")
	case mode&os.ModeNamedPipe != 0:
		fmt.Println("named pipe")
	}

}
func toHTTPError(err error) (msg string, httpStatus int) {
	if os.IsNotExist(err) {
		return "404 page not found", http.StatusNotFound
	}
	if os.IsPermission(err) {
		return "403 Forbidden", http.StatusForbidden
	}
	typ := fmt.Sprintf("%T", err)
	if typ != "" {
		return typ, http.StatusBadRequest
	}
	// Default:
	return "500 Internal Server Error", http.StatusInternalServerError
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
