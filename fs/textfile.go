package fs

import (
	"bytes"
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

func (f *File) AsTextfile() (t *Textfile, e error) {

	f.Type = "FT"
	t = &Textfile{File: f, Title: strings.TrimSuffix(f.Name, filepath.Ext(f.Name))}
	e = t.Load()
	return
}

type Textfile struct {
	*File
	Title   string
	Content []byte
}

func (f *Textfile) Load() (e error) {

	body, e := ioutil.ReadFile(f.location)
	if utf8.Valid(body) {
		f.Content = body
	} else {
		e = errors.New("Invalid UTF-8")
	}
	return
}

func (f *Textfile) StringContent() string {
	return string(f.Content)
}

func (f *Textfile) Save() error {
	return ioutil.WriteFile(f.location, f.Content, 0600)
}

func (t *Textfile) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//
	if r.Method == http.MethodPost {
		content := []byte(r.FormValue("content"))
		if len(content) > 0 && !bytes.Equal(content, t.Content) {
			if !utf8.Valid(content) {
				http.Error(w, "Invalid UTF-8", http.StatusBadRequest)
				return
			}
			t.Content = content
			if err := t.Save(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	if r.Method == http.MethodPut {
		body, _ := ioutil.ReadAll(r.Body)
		t.Content = body
		err := t.Save()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	t.Render(w, r)
}

func (t *Textfile) Render(w http.ResponseWriter, request *http.Request) {

	switch request.Header.Get("Content-type") {

	case "application/json":

		json, err := json.Marshal(t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(json)

	default:

		err := views.Render("textfile", w, t)
		if err != nil {
			fmt.Println("ERROR:", err)
		}
	}

}
