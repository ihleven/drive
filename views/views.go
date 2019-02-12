// Package views provides primitives for caching and rendering templates.
package views

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"path/filepath"
)

var viewMap map[string]*View
var Image, TextView *View
var RenderTextFile, RenderDir func(w http.ResponseWriter, data interface{}) error

func init() {
	viewMap = make(map[string]*View)
	Register("textfile", FuncMap, "./templates/textfile.html")
	//Register("image", FuncMap, "./vue/dist/file.html", "./templates/image.html", "./templates/layout/breadcrumbs.html")
	//Register("album", FuncMap, "./vue/dist/album.html")
	//Register("diary", FuncMap, "./vue/dist/diary.html")

	//RenderTextFile = RegisterRenderFunc(FuncMap, "./vue/dist/file.html", "./templates/textfile.html")
	//RenderDir = RegisterRenderFunc(FuncMap, "./vue/dist/directory.html", "./templates/layout/breadcrumbs.html")
	RenderDir = RegisterRenderFunc(FuncMap, "./templates/directory.html", "./templates/layout/breadcrumbs.html")

}

// A View is registered on startup and holds a template under a certain name to be rendered
// either directrly with (*view) Render or
// indirectly with Render and the lookup viewMap
type View struct {
	Name     string
	FuncMap  template.FuncMap
	Template *template.Template
}

func (v *View) Render(w http.ResponseWriter, data interface{}) (err error) {

	err = v.Template.Execute(w, data)
	//err = v.Template.ExecuteTemplate(w, "file", data)
	return
}

func Register(viewname string, funcMap template.FuncMap, files ...string) *View {
	name := filepath.Base(files[0])
	tmpl := template.Must(template.New(name).Funcs(funcMap).ParseFiles(files...))
	view := &View{Name: viewname, FuncMap: funcMap, Template: tmpl}
	viewMap[viewname] = view
	return view
}
func Render(viewname string, w http.ResponseWriter, data interface{}) error {
	if view, ok := viewMap[viewname]; ok {
		//fmt.Println("rendering view:", viewname, view.Template.Name(), view.Template)

		err := view.Template.Execute(w, data)
		//err := v.Template.ExecuteTemplate(w, "file", data)
		return err
	}
	return errors.New("not found")
}

// Fprint formats using the default formats for its operands and writes to w.
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func RegisterRenderFunc(funcMap template.FuncMap, files ...string) func(w http.ResponseWriter, data interface{}) error {
	name := filepath.Base(files[0])
	tmpl := template.Must(template.New(name).Funcs(funcMap).ParseFiles(files...))
	//fmt.Println("REgister", tmpl)

	return func(w http.ResponseWriter, data interface{}) error {
		//fmt.Println("rendering:", tmpl)

		err := tmpl.Execute(w, data)
		//fmt.Println(err)
		return err
	}
}

func SerializeJSON(w http.ResponseWriter, fh interface{}) {

	json, err := json.Marshal(fh)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(json)
}
