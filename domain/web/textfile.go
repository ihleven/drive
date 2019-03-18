package web

import (
	"bytes"
	"drive/domain"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/gomarkdown/markdown"
)

type TextFileController struct {
	File          *domain.File
	User          *domain.Account
	Title         string
	Content       []byte
	StringContent string
}

func (h TextFileController) Init(file *domain.File, usr *domain.Account, st domain.Storage) {
	h.File = file
	h.User = usr
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
	fmt.Println("TextFileController")
	title := strings.TrimSuffix(c.File.Name, filepath.Ext(c.File.Name))
	content, err := c.File.GetContent()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output := template.HTML(markdown.ToHTML(content, nil, nil))
	stringContent, _ := c.File.GetUTF8Content()
	m := map[string]interface{}{"user": c.User, "file": c.File, "content": stringContent, "Title": title, "Markdown": output}

	switch r.Header.Get("Content-type") {

	case "application/json":

		json, err := json.Marshal(c.File)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(json)
	default:
		fmt.Println("TextFileController")
		rnd.HTML(w, http.StatusOK, "file2", m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
