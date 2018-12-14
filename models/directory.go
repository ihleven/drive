package models

import (
	"drive/templates"
	"encoding/json"
	"errors"
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

func NewFiler(root, pathname string) (Filer, error) {

	basefile := NewFilePath(root, pathname)
	if basefile == nil {
		return nil, errors.New("404")
	}
	if (*basefile.Info).IsDir() {
		return NewDirectory(basefile)
	} else {
		return NewFile(basefile), nil
	}

}

type FilePath struct {
	Info     *os.FileInfo
	Abs      string // /home/ihle/alben/urlaube/Wien2013/IMG_123.jpg
	Path     string // /alben/urlaube/Wien2013/IMG_123.jpg
	Dirname  string // /alben/urlaube/Wien2013/
	Filename string // IMG_123.jpg
	Root     string // /home/ihle
}

func NewFilePath(root, pathname string) *FilePath {

	fullpath := path.Join(root, pathname)
	fmt.Println(fullpath)
	fileInfo, error := os.Stat(fullpath)
	if error != nil {
		if os.IsNotExist(error) {
			return nil
		}
		panic(fmt.Sprint("NewFiler -> os.Stat Error:", error))
	}
	dirname, filename := path.Split(fullpath)
	return &FilePath{Root: root, Path: pathname, Abs: fullpath, Dirname: dirname, Filename: filename, Info: &fileInfo}
}

type Directory struct {
	FilePath
	Parent    string
	Folders   []string
	Files     []*File
	IndexFile string
}

func NewDirectory(fp *FilePath) (*Directory, error) {
	d := &Directory{FilePath: *fp}
	d.scan()
	return d, nil
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
	fmt.Printf(" - scanning '%s'\n", d.Path)
	fileInfos, err := ioutil.ReadDir(d.Abs)
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
			fp := NewFilePath(d.FilePath.Root, path.Join(d.FilePath.Path, info.Name()))
			file := NewFile(fp)
			if err := file.Scan(); err == nil {
				if file.MIME.Type == "image" {
					d.Files = append(d.Files, file)
				} else {
					d.Files = append(d.Files, file)
				}
			}

		}
	}
	return nil
}
