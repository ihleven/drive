package views

import (
	"drive/config"
	"drive/fs"
	"drive/templates"
	"encoding/json"
	"fmt"
	"net/http"
)

type Renderer interface {
	render(w http.ResponseWriter, req *http.Request)
}

type View struct {
	template string
}

func PathHandler(w http.ResponseWriter, r *http.Request) {

	var tmpl string
	fd, _ := fs.NewPath(config.Root, r.URL.Path)
	if fd == nil {
		http.Error(w, fmt.Sprintf("Not found: %s", r.URL.Path), http.StatusNotFound)
		return
	}
	if fd.IsDir() {
		fd.List()
		tmpl = "directory.html"
	} else if fd.IsRegular() {
		if fd.MIME.Type == "text" {
			textfile := fd.Textfile()
			templates.RenderTemplate(w, "file.html", textfile)
		}
		tmpl = "file.html"
	}

	err := templates.RenderTemplate(w, tmpl, fd)
	if err != nil {
		fmt.Println(err)
	}

	return
	res, _ := json.Marshal(fd)
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
