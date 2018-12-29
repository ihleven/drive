package fs

import (
	"drive/views"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

type Directory struct {
	File
	//Parent    string
	Folders   []File
	Files     []File
	Children  []File
	IndexFile string
}

func (d *Directory) Filetype() string {
	return "directory"
}
func NewDirectory(fi *File) (*Directory, error) {
	d := &Directory{File: *fi}
	//d.scan()
	return d, nil
}

func (d *Directory) NewChildFromFileInfo(info os.FileInfo) *File {

	file := File{
		location: path.Join(d.location, info.Name()),
		Path:     path.Join(d.Path, info.Name()),
		Name:     info.Name(),
		Size:     info.Size(),
		Mode:     info.Mode(),
		ModTime:  info.ModTime()}
	if info.IsDir() {
		file.Type = "D"
		d.Folders = append(d.Folders, file)
	} else {
		file.GuessMIME()
		//child.ParseMIME()
		//child.MatchMIMEType()
		//child.DetectContentType()
		d.Files = append(d.Files, file)
	}
	return &file
}

func (d *Directory) List() error {

	entries, err := ioutil.ReadDir(d.location)
	if err != nil {
		log.Fatal(err)
		return err
	}
	for _, info := range entries {
		if info.Name()[0] == '.' {
			continue
		}
		d.NewChildFromFileInfo(info)

	}
	d.Children = append(d.Folders, d.Files...)
	return nil
}

func (d *Directory) Handle(w http.ResponseWriter, r *http.Request) {

	d.List()
	switch r.Method {
	case http.MethodGet:
		d.Render(w, r)
	case http.MethodPost:
		fmt.Fprintln(w, "POST dir:", d.Name)
	}
}

func (f *File) SerializeJSON(w http.ResponseWriter, fh FileHandler) {

	json, err := json.Marshal(fh)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

func (d *Directory) Render(w http.ResponseWriter, request *http.Request) {

	switch request.Header.Get("Accept") {

	case "application/json":

		d.SerializeJSON(w, d)

	default:

		err := views.RenderDir(w, d)
		if err != nil {
			panic(err)
		}
	}

}
