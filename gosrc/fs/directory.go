package fs

import (
	"drive/gosrc/views"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
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
func NewDirectory(file *File) (*Directory, error) {
	file.Type = "D"
	dir := &Directory{File: *file}
	dir.List()
	return dir, nil
}

func (d *Directory) List() error {
	fmt.Println("LIST")
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

func (d *Directory) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var options struct {
			CreateThumbnails bool
		}
		err := decoder.Decode(&options)
		if err != nil {
			http.Error(w, "Bad Request", 400)
			return
		}
		if options.CreateThumbnails {
			fmt.Println(filepath.Join(d.location, "thumbs"))
			os.MkdirAll(filepath.Join(d.location, "thumbs"), os.ModePerm)

			for i := 0; i < len(d.Files); i++ {
				file := d.Files[i]

				if file.MIME.Type == "image" {
					img, err := file.AsImage()
					if err != nil {
						fmt.Println(" - ERROR '%s'\n\n\n", err)
						continue
					}
					img.MakeThumbnail()
				}

			}

		}

	}
	d.Render(w, r)
}

func (d *Directory) Render(w http.ResponseWriter, request *http.Request) {

	switch request.Header.Get("Accept") {

	case "application/json":

		views.SerializeJSON(w, d)

	default:

		err := views.RenderDir(w, d)
		if err != nil {
			panic(err)
		}
	}

}
