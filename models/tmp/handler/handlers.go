package handler

import (
	"drive/file"
	"drive/file/storage"
	"drive/models/auth"
	"drive/session"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"strings"
)

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		user, err := auth.Authenticate(r.PostFormValue("username"), r.PostFormValue("password"))
		if HttpLogOnError(w, err, "") {
			return
		}
		fmt.Println("=>", user, r.PostFormValue("username"), r.PostFormValue("password"))

		session.SetSessionUser(r, w, user)
		if HttpLogOnError(w, err, "session error") {
			fmt.Println("saving error", err)
			return
		}
		if next, ok := r.URL.Query()["next"]; ok && len(next) >= 1 {
			http.Redirect(w, r, next[0], 302)
		}
		http.Redirect(w, r, "/", 302)
	}

	t, _ := template.ParseFiles("templates/login.html")
	t.Execute(w, nil)
}

func Logout(w http.ResponseWriter, r *http.Request) {
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

	path := strings.TrimPrefix(path.Clean(r.URL.Path), "/serve/")

	file, err := file.Open(path, os.O_RDONLY, authuser.Uid, authuser.Gid)
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
	http.ServeContent(w, r, file.Descriptor.Name(), file.ModTime(), file.Descriptor)
}

func Raw(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(path.Clean(r.URL.Path), "/serve/")

	fh, err := storage.DefaultStorage.Open(path)
	if err != nil {
		http.Error(w, err.Error(), 500)

		return
	}
	defer fh.File.Close()
	fmt.Println("fh", fh)
	http.ServeContent(w, r, fh.Name(), fh.ModTime(), fh.File)
}

func PathHandler(w http.ResponseWriter, r *http.Request) {

	usr, _ := session.GetSessionUser(r, w)

	info, err := file.Open(path.Clean(r.URL.Path), os.O_RDONLY, usr.Uid, usr.Gid)
	fmt.Println(info, r.URL.Path, os.O_RDONLY)
	if HttpLogOnError(w, err, "error opening file") {
		return
	}
	defer info.Close()

	f := file.FileFromInfo(info)

	f.Path = path.Clean(r.URL.Path)
	switch {
	case f.Mode.IsRegular():

		controller := FileController{File: f, User: usr}
		controller.ServeHTTP(w, r)

	case f.Mode.IsDir():

		controller := DirController{File: f, User: usr}
		controller.Render(w, r)

	case f.Mode&os.ModeSymlink != 0:
		fmt.Println("symbolic link")
	case f.Mode&os.ModeNamedPipe != 0:
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
func HttpLogOnError(w http.ResponseWriter, err error, message string) bool {
	if err == nil {
		return false
	}
	msg, code := toHTTPError(err)
	if message != "" {
		msg = fmt.Sprintf("%s => %s", msg, message)
	}
	http.Error(w, msg, code)
	return true
}
