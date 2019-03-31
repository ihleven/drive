package web

import (
	"drive/domain"
	"encoding/json"
	"html/template"
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

	// content := []byte(r.FormValue("content"))
	// fmt.Println("POST", content, len(content), !bytes.Equal(content, c.Content))
	// if len(content) > 0 && !bytes.Equal(content, c.Content) {
	// 	if !utf8.Valid(content) {
	// 		http.Error(w, "Invalid UTF-8", http.StatusBadRequest)
	// 		return
	// 	}
	// 	c.Content = content
	// 	if err := c.File.SetContent(content); err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// }

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
