package web

import (
	"drive/domain"
	"drive/views"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/thedevsaddam/renderer"
)

// https://github.com/thedevsaddam/renderer
var rnd *renderer.Render

func init() {
	fm := make([]template.FuncMap, 1)
	fm[0] = views.FuncMap
	fmt.Println(fm[0])
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
	Entries []*domain.File
	Folder  *domain.Folder
}

func (h DirHandler) Render(w http.ResponseWriter, r *http.Request) {

	folder, _ := domain.NewDirectory(h.File, h.User)
	//for _, _ := range folder.Entries {
	//fmt.Println(" ==> ", i, e)
	//}
	switch r.Header.Get("Accept") {

	case "application/json":

		views.SerializeJSON(w, h.Folder)

	default:
		m := map[string]interface{}{"user": h.User, "file": h.File, "dir": folder}
		err := views.RenderDir(w, m)
		if err != nil {
			fmt.Println("render error: ", err)
		}
	}

}

func (h DirHandler) Post(w http.ResponseWriter, r *http.Request) {

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

			// file.Mkdir(filepath.Join(d.File.Path, "thumbs"))

			// for i := 0; i < len(d.Directory.Entries); i++ {
			// 	file := d.Directory.Entries[i]

			// 	if file.MIME.Type == "image" {
			// 		//img, err := file.AsImage()

			// 		//img.MakeThumbnail()
			// 	}

			// }

		}

	}

}
