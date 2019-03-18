package web

import (
	"drive/domain"
	"drive/domain/usecase"
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

	opts := renderer.Options{
		FuncMap:          []template.FuncMap{views.FuncMap},
		ParseGlobPattern: "./templates/*.html",
	}

	rnd = renderer.New(opts)
}

type ViewSet interface {
	Init(*domain.File, *domain.Account, domain.Storage)
	Render(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
}

type FileHandler struct {
	file *domain.File
	usr  *domain.Account
}

func (h *FileHandler) Init(file *domain.File, usr *domain.Account, st domain.Storage) {
	fmt.Println("filehandler", file, usr)
	h.file = file
	h.usr = usr

}
func (h *FileHandler) Render(w http.ResponseWriter, r *http.Request) {
	m := map[string]interface{}{"user": h.usr, "file": h.file}
	fmt.Println("filehandler render", h.file)

	err := rnd.HTML(w, http.StatusOK, "file", m)
	if err != nil {
		fmt.Println("render error: ", err)
	}
}

func (h *FileHandler) Post(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hgfghfjh"))
}

type DirHandler struct {
	File    *domain.File
	User    *domain.Account
	Entries []*domain.File
	Folder  *domain.Folder
}

func (h *DirHandler) Init(file *domain.File, usr *domain.Account, storage domain.Storage) {
	h.File = file
	h.User = usr
	folder, _ := usecase.GetFolder(storage, file, usr)
	h.Folder = folder
}
func (h *DirHandler) Render(w http.ResponseWriter, r *http.Request) {
	//folder, _ := domain.NewDirectory(h.File, h.User)
	//for _, _ := range folder.Entries {
	//fmt.Println(" ==> ", i, e)
	//}
	switch r.Header.Get("Accept") {

	case "application/json":

		views.SerializeJSON(w, h.Folder)

	default:
		m := map[string]interface{}{"user": h.User, "file": h.File, "dir": h.Folder}
		fmt.Println("dirhandler", m)

		err := rnd.HTML(w, http.StatusOK, "directory", m)
		if err != nil {
			fmt.Println("render error: ", err)
		}
	}

}

func (h *DirHandler) Post(w http.ResponseWriter, r *http.Request) {

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
