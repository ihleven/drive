package web

import (
	"drive/domain"
	"drive/domain/usecase"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ImageController struct {
	File     *domain.File
	User     *domain.Account
	Image    *usecase.Image
	Siblings *domain.Siblings
}

func (i *ImageController) Init(f *domain.File, u *domain.Account, s domain.Storage) {
	i.File = f
	i.User = u

	image, err := usecase.NewImage(f, u)
	if err != nil {
		return
	}
	i.Image = image
}

func (i *ImageController) Render(w http.ResponseWriter, r *http.Request) {
	siblings, err := i.File.Siblings()
	if err != nil {
		fmt.Println("siblings error:", err)
	}
	i.Siblings = siblings

	if r.Method == http.MethodPut {
		body, _ := ioutil.ReadAll(r.Body)
		//i.update(body)
		json, _ := json.Marshal(body)
		w.Write(json)
	}
	if r.Method == http.MethodPost {
		i.Post(w, r)
	}

	switch r.Header.Get("Accept") {
	case "application/json":
		rnd.JSON(w, http.StatusOK, i)
	default:
		m := map[string]interface{}{"user": i.User, "file": i.File, "Image": i.Image, "Siblings": i.Siblings}

		err := rnd.HTML(w, http.StatusOK, "image", m)
		if err != nil {
			fmt.Println("render error: ", err)
		}
	}
}
func (i *ImageController) Post(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	i.Image.Title = r.FormValue("title")
	i.Image.Caption = r.FormValue("caption")
	i.Image.Cutline = r.FormValue("cutline")
	if err := i.Image.WriteMeta(i.User); err != nil {
		fmt.Println("POST:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("POST: no error")
	http.Redirect(w, r, i.File.Path, http.StatusFound)

}
