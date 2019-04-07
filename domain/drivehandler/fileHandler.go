package drivehandler

import (
	"drive/domain"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/pkg/errors"
)

type FileHandler struct {
	File    *domain.File
	User    *domain.Account
	Title   string
	Content string
}

func (h *FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	h.Title = strings.TrimSuffix(h.File.Name, filepath.Ext(h.File.Name))
	content, _ := h.File.GetUTF8Content()
	h.Content = content

	h.Respond(w, r.Header.Get("Accept"))
}

func (h *FileHandler) Respond(w http.ResponseWriter, format string) error {

	var err error

	switch format {
	case "application/json":
		err = rnd.JSON(w, http.StatusOK, h)
	default:
		err = h.Render(w, h)
	}
	if err != nil {
		return errors.Wrap(err, "render error")
	}
	return nil

}

func (h *FileHandler) Render(w http.ResponseWriter, m interface{}) error {

	err := rnd.HTML(w, http.StatusOK, "file", m)
	if err != nil {
		return errors.Wrap(err, "render error")
	}
	return nil
}

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

	file, handle, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	_ = handle.Header.Get("Content-Type")

	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = c.File.SetUTF8Content(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
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
