package web

import (
	"drive/domain"
	"drive/domain/usecase"
	"drive/session"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/h2non/filetype/types"
)

func DispatchStorage(storage domain.Storage) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		sessionUser, _ := session.GetSessionUser(r, w)

		file, err := usecase.GetFile(storage, path.Clean(r.URL.Path), sessionUser)
		if err != nil {
			ErrorResponder(w, "error opening file: "+err.Error(), 500)
			return
		}

		var handler = GetRegisteredHandler(&file.MIME)
		handler.Init(file, sessionUser, storage)
		switch r.Method {
		case "GET":
			handler.Render(w, r)
		case "POST":
			handler.Post(w, r)
		}

		switch {
		case file.Mode.IsRegular():

			//controller := FileController{File: file, User: usr}
			//controller.ServeHTTP(w, r)

		case file.Mode.IsDir():

			//controller := DirController{File: f, User: usr}
			//controller.Render(w, r)

		case file.Mode&os.ModeSymlink != 0:
			fmt.Println("symbolic link")
		case file.Mode&os.ModeNamedPipe != 0:
			fmt.Println("named pipe")
		}

	}
}

func GetRegisteredHandler(MIME *types.MIME) (vs ViewSet) {

	//m := file.GuessMIME()
	switch MIME.Type {
	case "directory":
		vs = &DirHandler{}
	case "text":
		vs = &TextFileController{}
	case "image":
		vs = &ImageController{}
		//return ImageHandler{file: file, usr: usr}
	default:
		vs = &FileHandler{}
	}
	return
}

func Serve(storage domain.Storage) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		authuser, err := session.GetSessionUser(r, w)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		cleanedPath := path.Clean(r.URL.Path)[6:] // strip "/serve"-prefix

		handle, err := usecase.GetReadHandle(storage, cleanedPath, authuser.Uid, authuser.Gid)
		if err != nil {
			ErrorResponder(w, err.Error(), 500)
			return
		}

		if handle.IsDir() {
			r.URL.Path = path.Join(r.URL.Path, "index.html")
			Serve(storage)(w, r)
			return
		}

		fd := handle.Descriptor(0)
		defer fd.Close()

		http.ServeContent(w, r, handle.Name(), handle.ModTime(), fd)
	}
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

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	usr, _ := session.GetSessionUser(r, w)

	rnd.HTML(w, http.StatusOK, "index", map[string]interface{}{"user": usr})
}

func assetHandler(prefix, location string) http.Handler {

	return http.StripPrefix(fmt.Sprintf("/%s/", prefix), http.FileServer(http.Dir(fmt.Sprintf("./%s", location))))
}
