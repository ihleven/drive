package drivehandler

import (
	"drive/domain"
	"drive/domain/usecase"
	"drive/web"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type DirHandler struct {
	File   *domain.File
	User   *domain.Account
	Folder *domain.Folder
}

func (h *DirHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	format := r.Header.Get("Accept")

	folder, err := usecase.GetFolder(h.File, h.User)
	if err != nil {
		web.HandleError(w, r, err)
	}
	h.Folder = folder

	switch r.Method {
	case http.MethodGet:
		err = h.Get(r)
	case http.MethodPost:
		err = h.Post(r)
		if err == nil {
			http.Redirect(w, r, h.File.Path, http.StatusFound)
		}
	}
	if err != nil {
		web.HandleError(w, r, err)
	}

	err = h.Respond(w, format)
	if err != nil {
		web.HandleError(w, r, err)
	}

}

func (h *DirHandler) Post(r *http.Request) error {

	decoder := json.NewDecoder(r.Body)
	var options struct {
		CreateThumbnails bool
	}
	err := decoder.Decode(&options)
	if err != nil {
		return err
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
	return nil
}
func (h *DirHandler) Get(r *http.Request) error {
	return nil
}
func (h *DirHandler) Respond(w http.ResponseWriter, format string) error {
	var err error
	m := map[string]interface{}{"User": h.User, "File": h.File, "Folder": h.Folder}

	switch format {
	case "application/json":
		err = rnd.JSON(w, http.StatusOK, h)
	default:
		err = h.Render(w, m)
	}
	if err != nil {
		return errors.Wrap(err, "render error")
	}
	return nil
}

func (h *DirHandler) Render(w http.ResponseWriter, m interface{}) error {

	err := rnd.HTML(w, http.StatusOK, "folder", m)
	if err != nil {
		return errors.Wrap(err, "render error")
	}
	return nil
}
