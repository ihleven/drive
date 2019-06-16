package templates

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/thedevsaddam/renderer"
)

// https://github.com/thedevsaddam/renderer
var rnd *renderer.Render

func init() {

	opts := renderer.Options{
		FuncMap:          []template.FuncMap{FuncMap},
		ParseGlobPattern: "./_static/*.html",
	}
	rnd = renderer.New(opts)
}

func SerializeJSON(w http.ResponseWriter, status int, v interface{}) error {
	return rnd.JSON(w, status, v)
}

func Render(w http.ResponseWriter, status int, name string, v interface{}) error {
	err := rnd.HTML(w, status, name, v)
	fmt.Println(err)
	return err
}

//var Render = Rnd.HTML
