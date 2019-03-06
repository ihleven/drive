package handler

import (
	"drive/file"
	"drive/views"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
)

type DirController struct {
	File      *file.File
	Directory *file.Folder
	//User      *auth.Account
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
