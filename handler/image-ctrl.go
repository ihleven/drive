package handler

import (
	"drive/auth"
	"drive/file"
	"drive/models"
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
