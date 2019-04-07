package templates

import (
	"html/template"

	"github.com/thedevsaddam/renderer"
)

// https://github.com/thedevsaddam/renderer
var Rnd *renderer.Render

func init() {

	opts := renderer.Options{
		FuncMap:          []template.FuncMap{FuncMap},
		ParseGlobPattern: "./_static/templates/*.html",
	}

	Rnd = renderer.New(opts)
}
