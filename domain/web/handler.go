package web

import (
	"drive/domain"
	"drive/domain/usecase"
	"drive/session"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func Dispatch(w http.ResponseWriter, r *http.Request) {

	usr, _ := session.GetSessionUser(r, w)

	file, err := usecase.GetFile("public", path.Clean(r.URL.Path))
	fmt.Println(file, err)
	if err != nil {
		ErrorResponder(w, "error openign file: "+err.Error(), 500)
		return
	}

	defer file.Close()

	var mimeHandler = GetRegisteredMIMEHandler(file, usr)
	switch r.Method {
	case "GET":

		mimeHandler.Render(w, r)

	case "POST":
		mimeHandler.Post(w, r)
	}

	fmt.Println(r.Method)

}
func GetRegisteredMIMEHandler(file *domain.File, usr *domain.Account) ViewSet {
	m := file.GuessMIME()
	switch m.Type {
	case "text":
		return &FileHandler{file, usr}
	case "image":
		return &DirHandler{file, usr}
		//return ImageHandler{file: file, usr: usr}
	default:
		return FileHandler{}
	}
}

func Serve(w http.ResponseWriter, r *http.Request) {

	authuser, err := session.GetSessionUser(r, w)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	path := path.Clean(r.URL.Path)[6:] // strip "/serve"-prefix

	handle, err := usecase.GetServeContentHandle("public", path, authuser.Uid, authuser.Gid)
	if err != nil {
		ErrorResponder(w, err.Error(), 500)
		return
	}
	defer handle.Close()

	http.ServeContent(w, r, handle.Name(), handle.ModTime(), handle.GetFile())
}

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		user, err := usecase.Authenticate(r.PostFormValue("username"), r.PostFormValue("password"))
		if err != nil {
			ErrorResponder(w, err.Error(), 500)
			return
		}
		fmt.Println("=>", user, r.PostFormValue("username"), r.PostFormValue("password"))

		session.SetSessionUser(r, w, user)
		if err != nil {
			ErrorResponder(w, "session error"+err.Error(), 500)
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

func assetHandler(prefix, location string) http.Handler {

	return http.StripPrefix(fmt.Sprintf("/%s/", prefix), http.FileServer(http.Dir(fmt.Sprintf("./%s", location))))
}

func AlbumHandler(w http.ResponseWriter, r *http.Request) {

	path, _ := filepath.Rel("/alben", path.Clean(r.URL.Path))
	fmt.Printf(" - scanning '%s'\n", "/"+path)

	file, err := storage.Open("/" + path)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), 500)
	}

	if file.IsDir() {
		dir, _ := fs.NewDirectory(file)
		album, _ := fs.NewAlbum(dir)
		album.Render(w, r)
		return
	}
	if file.IsRegular() {

		diary, _ := fs.NewDiary(file, storage)
		fmt.Println("DIARY", diary)
		diary.ServeHTTP(w, r)
	}

}
