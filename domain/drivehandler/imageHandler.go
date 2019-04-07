package drivehandler

import (
	"drive/domain"
	"drive/domain/usecase"
	"drive/web"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type DriveHandler interface {
	http.Handler
	//Init() error
	Get(r *http.Request) error
	Post(r *http.Request) error
	Respond(w http.ResponseWriter, format string) error
	//Render(w http.ResponseWriter)
}

type ImageHandler struct {
	File     *domain.File
	User     *domain.Account
	Image    *usecase.Image
	Siblings *domain.Siblings
}

func (h *ImageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	format := r.Header.Get("Accept")
	err := h.Init()
	if err != nil {
		web.HandleError(w, r, err)
	}

	switch r.Method {
	case http.MethodGet:
		err = h.Get(r)
	case http.MethodPost:
		err = h.Post(r)
		if err == nil {
			http.Redirect(w, r, h.File.Path, http.StatusFound)
		}

	case http.MethodPut:
		body, _ := ioutil.ReadAll(r.Body)
		//i.update(body)
		json, _ := json.Marshal(body)
		w.Write(json)
	}

	if err != nil {
		web.HandleError(w, r, err)
	}

	err = h.Respond(w, format)
	if err != nil {
		web.HandleError(w, r, err)
	}
}

func (h *ImageHandler) Init() error {
	image, err := usecase.NewImage(h.File, h.User)
	if err != nil {
		return errors.Wrap(err, "NewImage")
	}
	h.Image = image
	siblings, err := h.File.Siblings()
	if err != nil {
		return errors.Wrap(err, "siblings error")
	}
	h.Siblings = siblings
	return nil
}
func (h *ImageHandler) Get(r *http.Request) error {
	return nil
}
func (h *ImageHandler) Respond(w http.ResponseWriter, format string) error {
	var err error
	m := map[string]interface{}{"User": h.User, "File": h.File, "Image": h.Image, "Siblings": h.Siblings}

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

func (h *ImageHandler) Render(w http.ResponseWriter, m interface{}) error {

	err := rnd.HTML(w, http.StatusOK, "image", m)
	if err != nil {
		return errors.Wrap(err, "render error")
	}
	return nil
}
func (h *ImageHandler) Post(r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return errors.Wrap(err, "parse form")
	}
	h.Image.Title = r.FormValue("title")
	h.Image.Caption = r.FormValue("caption")
	h.Image.Cutline = r.FormValue("cutline")
	if err := h.Image.WriteMeta(h.User); err != nil {
		return errors.Wrap(err, "Image.WriteMeta")
	}
	fmt.Println("POST: no error")
	return nil
}
