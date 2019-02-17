package file

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

type TextfileController struct {
	*File
	Title   string
	Content []byte
}

func NewTextfileController(file *File) (*TextfileController, error) {

	c := &TextfileController{File: file, Title: strings.TrimSuffix(file.Name, filepath.Ext(file.Name))}
	return c, nil
}

func (c *TextfileController) Post(w http.ResponseWriter, r *http.Request) {

	content := r.FormValue("content")
	oldcontent, _ := c.File.GetTextContent()
	if len(content) > 0 && !(content == oldcontent) {
		if !utf8.Valid([]byte(content)) {
			http.Error(w, "Invalid UTF-8", http.StatusBadRequest)
			return
		}

		if err := c.File.SetContent(string(content)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (c *TextfileController) Put(w http.ResponseWriter, r *http.Request) {

}

func (c *TextfileController) Render(w http.ResponseWriter, request *http.Request) {

	//switch request.Header.Get("Content-type") {

	//case "application/json":

	json, err := json.Marshal(c.File)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(json)

	//default:

	//	err := views.Render("textfile", w, t)
	//	if err != nil {
	//		fmt.Println("ERROR:", err)
	//	}
	//}

}
