package web

import (
	"drive/auth"
	"drive/domain"
	"drive/views"
	"html/template"
	"net/http"

	"github.com/thedevsaddam/renderer"
)

// https://github.com/thedevsaddam/renderer
var rnd *renderer.Render

func init() {
	fm := make([]template.FuncMap, 1)
	fm[0] = views.FuncMap
	opts := renderer.Options{
		FuncMap:          fm,
		ParseGlobPattern: "./templates/*.html",
	}

	rnd = renderer.New(opts)
}

type ViewSet interface {
	Render(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
}

type FileHandler struct {
	file *domain.File
	usr  *auth.Account
}

func (h FileHandler) Render(w http.ResponseWriter, r *http.Request) {
	rnd.HTML(w, http.StatusOK, "file", nil)
}

func (h FileHandler) Post(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hgfghfjh"))
}
