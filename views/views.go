package views

import (
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
	TextView = Register("textfile", FuncMap, "static/file.html", "templates/textfile.html")
	Image = Register("image", FuncMap, "static/file.html", "templates/image.html")

	RenderTextFile = RegisterRenderFunc(FuncMap, "static/file.html", "templates/textfile.html")
	RenderDir = RegisterRenderFunc(FuncMap, "static/directory.html")

}

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
