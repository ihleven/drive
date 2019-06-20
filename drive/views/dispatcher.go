package drivehandler

import (
	"drive/domain"
	"drive/drive"
	"drive/errors"
	"drive/session"

	"net/http"
	"path"
)

func DispatchStorage(storage drive.Storage) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		sessionUser, _ := session.GetSessionUser(r, w)

		file, err := drive.GetFile(storage, r.URL.Path, sessionUser)
		if err != nil {
			errors.Error(w, r, err)
			return
		}

		if responder := GetActioneer(file, sessionUser); responder != nil {
			var err error
			switch r.Method {
			case http.MethodGet:

				err = responder.GetAction(r, w)
			case http.MethodDelete:
				err = responder.DeleteAction(r, w)
			case http.MethodPost:
				err = responder.PostAction(r, w)

			case http.MethodPut:
				err = responder.PutAction(r, w)
			}
			if err != nil {

				errors.Error(w, r, err)
			}
		}
	}
}

func GetActioneer(file *drive.File, sessionUser *domain.Account) Actioneer {
	//fmt.Println("Actioneer:", file)
	switch {
	case file.IsDir():
		return &DirActionResponder{File: file, User: sessionUser}
	case file.MIME.Type == "image":
		return &ImageView{File: file, User: sessionUser}
	case file.Mode.IsRegular():
		return &FileActionResponder{File: file, User: sessionUser}
		//case file.Mode&os.ModeSymlink != 0:
		//	fmt.Println("symbolic link")
		//case file.Mode&os.ModeNamedPipe != 0:
		//	fmt.Println("named pipe")
	}
	return nil
}

func Serve(storage drive.Storage) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		authuser, err := session.GetSessionUser(r, w)
		if err != nil {
			errors.Error(w, r, err)
			return
		}

		cleanedPath := storage.CleanServePath(r.URL.Path) //path.Clean(r.URL.Path)[6:] // strip "/serve"-prefix

		handle, err := drive.GetReadHandle(storage, cleanedPath, authuser.Uid, authuser.Gid)
		if err != nil {
			errors.Error(w, r, err)
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

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
