package web

import (
	"drive/domain"
	"drive/views"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/thedevsaddam/renderer"
)

// https://github.com/thedevsaddam/renderer
var rnd *renderer.Render

func init() {
	fm := make([]template.FuncMap, 1)
	fm[0] = views.FuncMap
	opts := renderer.Options{
		FuncMap:          fm,
		ParseGlobPattern: "./templates/*.html",
	}

	rnd = renderer.New(opts)
}

type ViewSet interface {
	Render(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
}

type FileHandler struct {
	file *domain.File
	usr  *domain.Account
}

func (h FileHandler) Render(w http.ResponseWriter, r *http.Request) {
	rnd.HTML(w, http.StatusOK, "file", nil)
}

func (h FileHandler) Post(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hgfghfjh"))
}

type DirHandler struct {
	File    *domain.File
	User    *domain.Account
	Entries []*file.File
	Folder  *models.Folder
}

func (c DirController) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	folder, _ := file.NewDirectory(c.File, c.User)
	c.Directory = folder
	switch r.Method {
	case http.MethodPost:

		fmt.Fprintf(w, "POST")
	}
	c.Render(w, r)

}

func (c DirController) Render(w http.ResponseWriter, r *http.Request) {

	switch r.Header.Get("Accept") {

	case "application/json":

		views.SerializeJSON(w, c.Directory)

	default:
		m := map[string]interface{}{"user": c.User, "file": c.File, "dir": c.Directory}
		err := views.RenderDir(w, m)
		if err != nil {
			panic(err)
		}
	}

}

func (d *DirController) Action(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var options struct {
			CreateThumbnails bool
		}
		err := decoder.Decode(&options)
		if err != nil {
			http.Error(w, "Bad Request", 400)
			return
		}
		if options.CreateThumbnails {

			file.Mkdir(filepath.Join(d.File.Path, "thumbs"))

			for i := 0; i < len(d.Directory.Entries); i++ {
				file := d.Directory.Entries[i]

				if file.MIME.Type == "image" {
					//img, err := file.AsImage()

					//img.MakeThumbnail()
				}

			}

		}

	}
	d.Render(w, r)
}
