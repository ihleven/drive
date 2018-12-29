package main

import (
	"drive/config"
	"drive/fs"
	"fmt"
	"log"
	"mime"
	"net/http"
	"path"
	"regexp"
)

var storage = &fs.FileSystemStorage{}

func main() {
	mime.AddExtensionType(".py", "text/python")
	mime.AddExtensionType(".go", "text/golang")
	mime.AddExtensionType(".json", "text/json")
	//dbf()
	config.ParseFlags()
	//templates.Init()

	http.Handle("/serve/", http.StripPrefix("/serve/", http.FileServer(http.Dir(config.Root))))
	//http.Handle("/dist/", http.StripPrefix("/dist/", http.FileServer(http.Dir("frontend/dist"))))
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
	http.Handle("/hello/", http.StripPrefix("/hello/", http.HandlerFunc(sayhelloName)))
	http.HandleFunc("/hallo/", sayhelloName)
	//http.Handle("/alben/", models.AlbumController{})
	//http.HandleFunc("/", views.PathHandler)
	storage.Location = config.Root
	http.HandleFunc("/alben/", AlbumHandler)
	http.HandleFunc("/", storageContorller) //http.StripPrefix("/drive", mux))
	//router := Router{}
	//
	http.ListenAndServe(config.Address.String(), nil)
	//http.ListenAndServe(config.Address.String(), http.HandlerFunc(pathRequestHandler))
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Hello %s", r.URL.Path)

}

func storageContorller(w http.ResponseWriter, r *http.Request) {

	route := path.Clean(r.URL.Path)

	exists, _ := storage.Exists(route)
	if !exists {
		http.Error(w, fmt.Sprintf("Not found: %s", r.URL.Path), http.StatusNotFound)
		// http.NotFound(w, request)
		return
	}
	//contentType := r.Header.Get("Content-type")

	// assume text/html
	fd, _ := storage.GetSpecific(route)

	switch t := fd.(type) {
	case *fs.Directory:
		//handleDir(w, r, fd.(*fs.Directory))
		fd.Handle(w, r)
	case *fs.File:
		handleFile(w, r, fd.(*fs.File))
	default:
		fmt.Printf("Don't know type %T\n", t)
	}

}

func handleFile(w http.ResponseWriter, r *http.Request, f *fs.File) {

	if f.MIME.Type == "image" {
		image, _ := f.AsImage()
		image.ServeHTTP(w, r)

	}
	if f.MIME.Type == "text" {

		textfile, _ := f.AsTextfile()
		textfile.ServeHTTP(w, r)
	}

}

// wiki.go
func logit(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Before")
		h.ServeHTTP(w, r) // call original
		log.Println("After")
	})
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}
func main_wiki() {
	//http.HandleFunc("/edit/", makeHandler(editHandler))
	//http.HandleFunc("/save/", makeHandler(saveHandler))
	//http.HandleFunc("/", makeHandler(viewHandler))

}
