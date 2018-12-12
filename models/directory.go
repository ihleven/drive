package models

import (
	"drive/templates"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

type Filer interface {
	Render(w http.ResponseWriter, r *http.Request)
	//ContentType() string
}

func NewFiler(root, route string) (Filer, error) {

	absPath := path.Join(root, route)

	fileInfo, error := os.Stat(absPath)
	if error != nil {
		fmt.Println("NewFiler -> os.Stat Error:", error)
		return nil, error
	}
	basefile := basefile{FileInfo: fileInfo, fullpath: absPath}

	if fileInfo.IsDir() {
		d := &Directory{basefile: basefile}
		d.scan()
		return d, nil
	}
	return nil, nil
}

type basefile struct {
	os.FileInfo
	fullpath string
}

type Directory struct {
	basefile
	Name   string
	Route  string
	Parent string

	Folders   []string
	Files     []*File
	IndexFile string
}

func (d *Directory) Render(w http.ResponseWriter, req *http.Request) {

	if contentType := req.Header.Get("Content-type"); contentType == "application/json" {
		json.NewEncoder(w).Encode(d)
	} else {
		err := templates.Templates.ExecuteTemplate(w, "directory.html", d)
		if err != nil {
			fmt.Println(err)
		}
	}
}
func (d *Directory) scan() error {
	fmt.Printf(" - scanning '%s'\n", d.Route)
	fileInfos, err := ioutil.ReadDir(d.fullpath)
	if err != nil {
		log.Fatal(err)
		return err
	}
	for _, info := range fileInfos {
		if info.Name()[0] == '.' {
			continue
		}
		mode := info.Mode()

		if mode.IsDir() {
			d.Folders = append(d.Folders, info.Name())

		} else if mode.IsRegular() {

			fd := &File{Path: path.Join(d.fullpath, info.Name()), Name: info.Name()}
			if err := fd.Scan(); err == nil {
				if fd.MIME.Type == "image" {
					d.Files = append(d.Files, fd)
				} else {
					d.Files = append(d.Files, fd)
				}
			}

		}
	}
	return nil
}
