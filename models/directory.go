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

func NewFiler(root, path string) (Filer, error) {

	absPath := path.Join(root, route)
	dirname, filename := path.Split(path)

	fileInfo, error := os.Stat(absPath)
	if error != nil {
		if os.IsNotExist(error) {
			return nil, nil
		}
		Fatal("NewFiler -> os.Stat Error:", error)
	}

	basefile := BaseFile{
		Abs:      absPath,
		Path:     path,
		Dirname:  dirname,
		Filename: filename}

	if fileInfo.IsDir() {
		return NewDirectory(basefile)
	} else {
		return NewFile(basefile)
	}
	return nil, nil
}

type Directory struct {
	FilePath
	Parent string

	Folders   []string
	Files     []File
	IndexFile string
}

func (d *Directory) Render(w http.ResponseWriter, req *http.Request) {

	if contentType := req.Header.Get("Content-type"); contentType == "application/json" {
		json.NewEncoder(w).Encode(d)
	} else {
		fmt.Println("asdfasdfasdf")
		templates.RenderTemplate(w, "directory.html", d)
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
			if info.Name() == "index.html" {
				d.IndexFile = info.Name()
			}
			fd := File{Path: path.Join(d.fullpath, info.Name()), Name: info.Name()}
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
