package main

import (
	"drive/auth"
	"drive/file"
	"drive/fs"
	"drive/session"
	"drive/views"
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

func Serve(w http.ResponseWriter, r *http.Request) {

	authuser, err := session.AuthUser(r, w)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	storage := &file.FileSystemStorage{Root: "/Users/mi/go"}
	filename := strings.TrimPrefix(path.Clean(r.URL.Path), "/serve/")

	file, err := storage.GetFile(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()

	if r, _, _ := file.Permissions(authuser.Uid, authuser.Gid); !r {
		http.Error(w, fmt.Sprintf("HTTP 403 Forbidden: user '%s' does not have permission to access '%s'", authuser.Username, filename), 403)
		return
	}
	if file.IsDir() {
		r.URL.Path = fmt.Sprintf("/serve/%s/index.html", filename)
		Serve(w, r)
		return
	}
	http.ServeContent(w, r, file.Name(), file.ModTime(), file.Descriptor)
}

func PathHandler(w http.ResponseWriter, r *http.Request) {

	usr, err := session.GetSessionUser(r, w)

	f, err := file.NewFile(path.Clean(r.URL.Path), usr)
	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return
	}
	defer f.Close()

	mode := f.Mode
	switch {
	case mode.IsRegular():
		fmt.Println("regular file")
		controller := FileController{f, usr}
		controller.Render(w, r)

	case mode.IsDir():
		fmt.Println("directory")
		controller := DirController{f, usr}
		controller.Render(w, r)
		//folder, _ := file.NewDirectory(f)
		//Dir(w, r, folder, usr)
	case mode&os.ModeSymlink != 0:
		fmt.Println("symbolic link")
	case mode&os.ModeNamedPipe != 0:
		fmt.Println("named pipe")
	}

}

type FileController struct {
	File *file.File
	User *auth.User
}

func (c FileController) Render(w http.ResponseWriter, r *http.Request) {
	content, err := c.File.GetTextContent()
	m := map[string]interface{}{"user": c.User, "file": c.File, "content": content}
	err = views.RenderFile(w, m)
	if err != nil {
		panic(err)
	}
}

type DirController struct {
	File *file.File
	User *auth.User
}

func (c DirController) Render(w http.ResponseWriter, r *http.Request) {

	folder, _ := file.NewDirectory(c.File)

	switch r.Method {
	case http.MethodPost:

		fmt.Fprintf(w, "POST")
	}
	// d.Render(w, r)
	switch r.Header.Get("Accept") {

	case "application/json":

		views.SerializeJSON(w, folder)

	default:
		m := map[string]interface{}{"user": c.User, "file": c.File, "dir": folder}
		err := views.RenderDir(w, m)
		if err != nil {
			panic(err)
		}
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
