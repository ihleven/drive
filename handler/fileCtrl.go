package handler

import (
	"drive/auth"
	"drive/file"
	"drive/views"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

type FileController struct {
	File    *file.File
	User    *auth.Account
	Content []byte
}

func (c FileController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	content, err := c.File.GetContent()
	fmt.Println("GetTextContent", content, err)

	m := map[string]interface{}{"user": c.User, "file": c.File, "content": content}
	err = views.RenderFile(w, m)
	if err != nil {
		panic(err)
	}
}

type TextFileController struct {
	File    *file.File
	User    *auth.Account
	Title   string
	Content string
}

func (t TextFileController) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	t.Title = strings.TrimSuffix(t.File.Name, filepath.Ext(t.File.Name))
	content, err := t.File.GetUTF8Content()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Content = content

	if r.Method == http.MethodPost {
		content := r.FormValue("content")
		if len(content) > 0 && content != t.Content {
			if !utf8.Valid([]byte(content)) {
				http.Error(w, "Invalid UTF-8", http.StatusBadRequest)
				return
			}
			t.Content = content
			if err := t.File.SetContent(content); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	t.Render(w, r)
}
func (c TextFileController) Render(w http.ResponseWriter, r *http.Request) {

	m := map[string]interface{}{"user": c.User, "file": c.File, "content": c.Content, "Title": c.Title}

	switch r.Header.Get("Content-type") {

	case "application/json":

		json, err := json.Marshal(c.File)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(json)
	default:
		err := views.RenderFile(w, m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
