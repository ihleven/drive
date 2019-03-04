package web

package handler

import (
	"drive/auth"
	"drive/file"
	"drive/file/storage"
	"drive/session"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"strings"
)

func Serve(w http.ResponseWriter, r *http.Request) {

	authuser, err := session.GetSessionUser(r, w)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	path := path.Clean(r.URL.Path)[6:] // strip "/serve"-prefix

	file, err := usecase.GetServeContentHandle(path, authuser.Uid, authuser.Gid)
	if err != nil {
		responder.ErrorPage(w, err.Error(), 500)
		return
	}
	defer file.Close()

	
	http.ServeContent(w, r, file.Name, file.MTime, file.Descriptor)
}
