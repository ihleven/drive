package handler

import (
	"drive/file"
	"drive/models"
	"drive/models/auth"
	"drive/views"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ImageController struct {
	File     *file.File
	User     *auth.Account
	Image    *models.Image
	Siblings *file.Siblings
}

func (i ImageController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	image, _ := models.NewImage(i.File)
	i.Image = image

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
		views.SerializeJSON(w, i)
	default:
		m := map[string]interface{}{"user": i.User, "file": i.File, "Image": i.Image, "Siblings": i.Siblings}

		err := views.Render("image", w, m)
		if err != nil {
			fmt.Println(err)
		}
	}
}
func (i ImageController) Post(w http.ResponseWriter, r *http.Request) {
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
