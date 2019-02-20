package controller

import (
	"bytes"
	"drive/auth"
	"drive/file"
	"drive/views"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

type FileController struct {
	File *file.File
	User *auth.User
}

func (c FileController) Render(w http.ResponseWriter, r *http.Request) {
	content, err := c.File.GetTextContent()
	fmt.Println("GetTextContent", content, err)

	m := map[string]interface{}{"user": c.User, "file": c.File, "content": content}
	err = views.RenderFile(w, m)
	if err != nil {
		panic(err)
	}
}

type TextFileController struct {
	File    *file.File
	User    *auth.User
	Title   string
	Content []byte
}

func (tf *TextFileController) Load() (e error) {

	body, e := tf.File.GetContent()
	if utf8.Valid(body) {
		tf.Content = body
	} else {
		e = errors.New("Invalid UTF-8")
	}
	return
}
func (c TextFileController) Render(w http.ResponseWriter, r *http.Request) {
	content, err := c.File.GetTextContent()
	fmt.Println("GetTextContent", content, err)

	m := map[string]interface{}{"user": c.User, "file": c.File, "content": content}
	err = views.RenderFile(w, m)
	if err != nil {
		panic(err)
	}

	switch r.Header.Get("Content-type") {

	case "application/json":

		json, err := json.Marshal(c.File)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(json)
	}
}

func (t *TextFileController) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	t.Title = strings.TrimSuffix(t.File.Name, filepath.Ext(t.File.Name))
	t.File.Type = "FT"

	if r.Method == http.MethodPost {
		content := []byte(r.FormValue("content"))
		if len(content) > 0 && !bytes.Equal(content, t.Content) {
			if !utf8.Valid(content) {
				http.Error(w, "Invalid UTF-8", http.StatusBadRequest)
				return
			}
			t.Content = content
			//if err := t.Save(); err != nil {
			//	http.Error(w, err.Error(), http.StatusInternalServerError)
			//	return
			//}
		}
	}
	if r.Method == http.MethodPut {
		body, _ := ioutil.ReadAll(r.Body)
		t.Content = body
		//err := t.Save()
		//if err != nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//	return
		//}
	}

	t.Render(w, r)
}

type DirController struct {
	File      *file.File
	Directory *file.Folder
	User      *auth.User
}

func (c DirController) Render(w http.ResponseWriter, r *http.Request) {

	folder, _ := file.NewDirectory(c.File, c.User)
	c.Directory = folder
	switch r.Method {
	case http.MethodPost:

		fmt.Fprintf(w, "POST")
	}
	// d.Render(w, r)
	switch r.Header.Get("Accept") {

	case "application/json":

		views.SerializeJSON(w, folder)

	default:
		m := map[string]interface{}{"user": c.User, "file": c.File, "dir": folder}
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
