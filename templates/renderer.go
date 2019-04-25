package templates

import (
	"html/template"
	"net/http"

	"github.com/thedevsaddam/renderer"
)

// https://github.com/thedevsaddam/renderer
var rnd *renderer.Render

func init() {

	opts := renderer.Options{
		FuncMap:          []template.FuncMap{FuncMap},
		ParseGlobPattern: "./_static/templates/*.html",
	}

	rnd = renderer.New(opts)
}

func SerializeJSON(w http.ResponseWriter, status int, v interface{}) error {
	return rnd.JSON(w, status, v)
}

func Render(w http.ResponseWriter, status int, name string, v interface{}) error {
	return rnd.HTML(w, status, name, v)
}

//var Render = Rnd.HTML
