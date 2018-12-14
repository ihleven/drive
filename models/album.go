package models

import (
	"drive/templates"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

var Root_folder string

func ValidPath(url string) bool {

	fp := path.Join(Root_folder, url)

	info, error := os.Stat(fp)
	if error != nil {
		fmt.Println("ValidPAth?", url, error)
		return false
	}
	fmt.Println("ValidPAth?", url, !info.IsDir())
	return !info.IsDir()
}

func PathHandler(w http.ResponseWriter, req *http.Request, path string) {
	fmt.Fprintf(w, "Path: %s", path)

}

type Album struct {
	Name        string
	Dirpath     string
	AlbumFile   string
	Title       string
	Description string
	Keywords    []string
	Files       []*File
	Images      []*File
}

func (a *Album) scan() error {
	fmt.Printf(" - scanning '%s'\n", a.Dirpath)
	fileInfos, err := ioutil.ReadDir(a.Dirpath)
	if err != nil {
		log.Fatal(err)
		return err
	}
	for _, info := range fileInfos {
		if info.Name()[0] == '.' {
			continue
		}
		if info.Name() == "album.html" {
			//handleFile(path.Join(f.Name(), "index.html"), w, req)
			a.AlbumFile = info.Name()
		}
		mode := info.Mode()

		if mode.IsDir() {
			//d.Dirs = append(d.Dirs, val.Name())

		} else if mode.IsRegular() {

		}
	}
	return nil
}
func (a *Album) render(w http.ResponseWriter, req *http.Request) error {

	contentType := req.Header.Get("Content-type")

	if contentType == "application/json" {
		//context.MarschalJSON(w, req)
	} else {

		fmt.Print("template")
		//err = tpl.Execute(w, a)
		err := templates.RenderTemplate(w, "album.html", a)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func HandleAlbum(w http.ResponseWriter, req *http.Request, name string) {

	//fmt.Fprintf(w, "Album: %s", name)

	dirpath := path.Join(Root_folder, name)
	album := &Album{Name: name, Dirpath: dirpath} //Dirs: make([]string, 0), Files: make([]string, 0)}

	if err := album.scan(); err != nil {
		http.Error(w, err.Error(), 500)
	}
	album.Title = fmt.Sprintf("%s | alben", name)

	album.render(w, req)
	return

}
