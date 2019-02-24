package app

import (
	"drive/file"
	"drive/handler"
	"drive/session"
	"fmt"
	"net/http"
	"os"
	"path"
)

type PathRouter struct {
}

var pathRouter PathRouter

func (pr PathRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	usr, _ := session.GetSessionUser(r, w)

	info, err := file.Open(path.Clean(r.URL.Path), os.O_RDONLY, usr.Uid, usr.Gid)
	if handler.HttpLogOnError(w, err, "error opening file") {
		return
	}
	defer info.Close()

	f := file.FileFromInfo(info)
	f.Path = path.Clean(r.URL.Path)

	if f.Mode.IsDir() {
		controller := handler.DirController{File: f, User: usr}
		controller.ServeHTTP(w, r)
		return
	}
	fmt.Println(f.MIME)
	var controller http.Handler
	switch f.MIME.Type {
	case "text":
		controller = handler.TextFileController{File: f, User: usr}

	case "image":
		controller = handler.ImageController{File: f, User: usr}

	default:
		controller = handler.FileController{File: f, User: usr}
	}
	controller.ServeHTTP(w, r)
	//case f.Mode&os.ModeSymlink != 0:
	//	fmt.Println("symbolic link")
	//case f.Mode&os.ModeNamedPipe != 0:
	//	fmt.Println("named pipe")
	//}
}
