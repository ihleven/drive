package drivehandler

import (
	"drive/domain"
	"drive/domain/storage"
	"drive/domain/usecase"
	"drive/session"
	"drive/web"
	"fmt"
	"net/http"
	"os"
	"path"
)

func Setup() {

	web.RegisterFunc("/login", Login)
	web.RegisterFunc("/logout", Logout)
	web.RegisterFunc("/serve/home/", Serve(storage.Get("home")))
	web.RegisterFunc("/serve/", Serve(storage.Get("public")))
	web.RegisterFunc("/public/", DispatchStorage(storage.Get("public")))
	web.RegisterFunc("/home/", DispatchStorage(storage.Get("home")))
}

func DispatchStorage(storage domain.Storage) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		sessionUser, _ := session.GetSessionUser(r, w)

		file, err := usecase.GetFile(storage, path.Clean(r.URL.Path), sessionUser)
		if err != nil {
			web.ErrorResponder(w, "error opening file: "+err.Error(), 500)
			return
		}

		var handler DriveHandler
		switch {
		case file.IsDir():
			handler = &DirHandler{File: file, User: sessionUser}
		case file.MIME.Type == "image":
			handler = &ImageHandler{File: file, User: sessionUser}
		case file.Mode.IsRegular():
		case file.Mode&os.ModeSymlink != 0:
			fmt.Println("symbolic link")
		case file.Mode&os.ModeNamedPipe != 0:
			fmt.Println("named pipe")
		default:
			//handler = FileHandler{File: file, User: sessionUser}
		}

		handler.ServeHTTP(w, r)
	}
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
			web.ErrorResponder(w, err.Error(), 500)
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
