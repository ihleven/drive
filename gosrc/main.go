package main

import (
	"drive/gosrc/config"
	"drive/gosrc/fs"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

func main() {

	//dbf()
	config.ParseFlags()
	//dbConf := config.GetDatabaseConfiguration("mssql")
	//InitStore(dbConf)

	http.Handle("/serve/", http.StripPrefix("/serve/", http.FileServer(http.Dir(config.Root))))
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("../vue/dist/assets"))))
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("../static"))))
	http.Handle("/hello/", http.StripPrefix("/hello/", http.HandlerFunc(sayhelloName)))
	http.HandleFunc("/hallo/", sayhelloName)
	//http.Handle("/alben/", models.AlbumController{})
	//http.HandleFunc("/", views.PathHandler)
	storage.Location = config.Root
	http.HandleFunc("/alben/", AlbumHandler)
	//http.HandleFunc("/", storageContorller) //http.StripPrefix("/drive", mux))
	http.Handle("/files/", http.StripPrefix("/files", storage)) //http.StripPrefix("/drive", mux))
	http.Handle("/", storage)                                   //http.StripPrefix("/drive", mux))
	//router := Router{}
	//

	http.HandleFunc("/api/accommodations", GetAccommodationsHandler)

	http.ListenAndServe(config.Address.String(), nil)
	//http.ListenAndServe(config.Address.String(), http.HandlerFunc(pathRequestHandler))
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Hello %s", r.URL.Path)

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

func AlbumHandler(w http.ResponseWriter, r *http.Request) {

	path, _ := filepath.Rel("/alben", path.Clean(r.URL.Path))
	fmt.Printf(" - scanning '%s'\n", "/"+path)

	file, err := storage.Open("/" + path)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), 500)
	}

	if file.IsDir() {
		dir, _ := fs.NewDirectory(file)
		album, _ := fs.NewAlbum(dir)
		album.Render(w, r)
		return
	}
	if file.IsRegular() {

		diary, _ := fs.NewDiary(file, storage)
		fmt.Println("DIARY", diary)
		diary.ServeHTTP(w, r)
	}

}
