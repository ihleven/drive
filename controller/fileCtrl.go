package controller

import (
	"drive/auth"
	"drive/file"
	"drive/views"
	"fmt"
	"net/http"
)

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

	folder, _ := file.NewDirectory(c.File, c.User)

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
