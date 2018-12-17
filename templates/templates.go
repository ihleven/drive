package templates

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var templates map[string]*template.Template

// https://blog.questionable.services/article/approximating-html-template-inheritance/
// https://github.com/asit-dhal/golang-template-layout

func bytes(size int64) string {
	ext := []string{"bytes", "kb", "mb", "gb"}
	i := 0
	for ; size > 1024; i++ {
		size = size / 1024
	}
	return fmt.Sprintf("%d %s", size, ext[i])
}

var funcMap = template.FuncMap{
	"bytes": bytes,
}

func Init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	//templatesDir := "templates"
	layouts, err := filepath.Glob("templates/*i*.html")
	if err != nil {
		log.Fatal(err)
	}

	includes, err := filepath.Glob("templates/layout.html")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(layouts, includes)
	// Generate our templates map from our layouts/ and includes/ directories
	for _, layout := range layouts {
		fmt.Println(filepath.Base(layout))
		//files := append(includes, layout)
		templates[filepath.Base(layout)] = template.Must(template.New(filepath.Base(layout)).Funcs(funcMap).ParseFiles("templates/layout.html", layout))
	}

}

// renderTemplate is a wrapper around template.ExecuteTemplate.
func RenderTemplate(w http.ResponseWriter, name string, data interface{}) error { //data map[string]interface{}) error {
	//fmt.Println(templates, data)
	// Ensure the template exists in the map.
	tmpl, ok := templates[name]
	if !ok {
		return fmt.Errorf("The template %s does not exist.", name)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tmpl.ExecuteTemplate(w, "layout", data)
	return err
}
