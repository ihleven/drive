package web

import (
	"drive/domain"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
)

type TextFileController struct {
	File          *domain.File
	User          *domain.Account
	Title         string
	Content       string
	StringContent template.HTML
	Markdown      string
}

func (h *TextFileController) Init(file *domain.File, usr *domain.Account, st domain.Storage) {
	h.File = file
	h.User = usr
}
func (c *TextFileController) Post(w http.ResponseWriter, r *http.Request) {

	if !c.File.Permissions.Write {
		e := errors.New("missing write permission").Error()
		http.Error(w, e, http.StatusForbidden)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	// mimeType := handle.Header.Get("Content-Type")

	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(string(data))

	err = c.File.SetContent(data)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
}

func (c *TextFileController) Render(w http.ResponseWriter, r *http.Request) {

	stringContent, _ := c.File.GetUTF8Content()
	content, _ := c.File.GetContent()

	c.Title = strings.TrimSuffix(c.File.Name, filepath.Ext(c.File.Name))
	c.StringContent = template.HTML(stringContent)
	c.Content = stringContent

	_ = template.HTML(markdown.ToHTML(content, nil, nil))

	_ = map[string]interface{}{
		"User":     c.User,
		"File":     c.File,
		"Content":  template.JS(stringContent),
		"Title":    c.Title,
		"Markdown": "",
	}

	switch r.Header.Get("Content-type") {

	case "application/json":

		json, err := json.Marshal(c.File)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(json)

	default:

		err := rnd.HTML(w, http.StatusOK, "textfile", c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
