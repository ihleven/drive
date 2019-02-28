package handler

import (
	"bytes"
	"drive/auth"
	"drive/file"
	"drive/views"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/gomarkdown/markdown"
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
	File          *file.File
	User          *auth.Account
	Title         string
	Content       []byte
	StringContent string
}

func (t TextFileController) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	t.Title = strings.TrimSuffix(t.File.Name, filepath.Ext(t.File.Name))
	content, err := t.File.GetContent()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Content = content

	if r.Method == http.MethodPost {
		t.Post(w, r)
	}

	t.Render(w, r)
}

func (c *TextFileController) PostAlt(w http.ResponseWriter, r *http.Request) {

	content := r.FormValue("content")
	oldcontent, _ := c.File.GetUTF8Content()
	if len(content) > 0 && !(content == oldcontent) {
		if !utf8.Valid([]byte(content)) {
			http.Error(w, "Invalid UTF-8", http.StatusBadRequest)
			return
		}

		if err := c.File.SetContent([]byte(content)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (c TextFileController) Post(w http.ResponseWriter, r *http.Request) {

	content := []byte(r.FormValue("content"))
	fmt.Println("POST", content, len(content), !bytes.Equal(content, c.Content))
	if len(content) > 0 && !bytes.Equal(content, c.Content) {
		if !utf8.Valid(content) {
			http.Error(w, "Invalid UTF-8", http.StatusBadRequest)
			return
		}
		c.Content = content
		if err := c.File.SetContent(content); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}

func (c TextFileController) Render(w http.ResponseWriter, r *http.Request) {

	output := template.HTML(markdown.ToHTML(c.Content, nil, nil))
	stringContent, _ := c.File.GetUTF8Content()
	m := map[string]interface{}{"user": c.User, "file": c.File, "content": stringContent, "Title": c.Title, "Markdown": output}

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
