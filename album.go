package main

import (
	"drive/fs"
	"encoding/json"
	"fmt"
	"html/template"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"path"
	"path/filepath"
)

type Album struct {
	*fs.Directory
	AlbumFile   string
	Title       string
	Description string
	Keywords    []string
	Images      []fs.Image
}

func AlbumHandler(w http.ResponseWriter, r *http.Request) {

	path, _ := filepath.Rel("/alben", path.Clean(r.URL.Path))
	fmt.Printf(" - scanning '%s'\n", path)

	dir, err := storage.GetDirectory(path)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	album := Album{Directory: dir}
	for _, file := range dir.Files {

		if file.Name == "album.html" {
			//handleFile(path.Join(f.Name(), "index.html"), w, req)
			album.AlbumFile = file.Name
		}
		if file.MIME.Type == "image" {
			img, err := file.AsImage()
			if err != nil {
				fmt.Println(" - ERROR '%s'\n\n\n", err)
			} else {
				album.Images = append(album.Images, *img)
			}
		}

	}

	album.Title = fmt.Sprintf("%s | alben", dir.Name)

	album.render(w, r)
}

func (a *Album) render(w http.ResponseWriter, req *http.Request) error {

	contentType := req.Header.Get("Content-type")

	if contentType == "application/json" {

		json, _ := json.Marshal(a)
		w.Write(json)
	} else {

		t := template.Must(template.New("index.html").ParseFiles("static/index.html"))
		err := t.Execute(w, a)

		if err != nil {
			panic(err)
		}
	}
	return nil
}
