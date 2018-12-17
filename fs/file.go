package fs

import (
	"drive/templates"
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/h2non/filetype"
)

func (p *Path) Textfile() *Textfile {
	file := &Textfile{Path: p}
	file.Load()
	return file
}

type Textfile struct {
	*Path
	Title   string
	Content []byte
}

func NewFile(fp *Path) *Textfile {

	f := &Textfile{Path: fp}
	f.Load()
	return f
}

func (f *Textfile) Render(w http.ResponseWriter, req *http.Request) {

	if req.Method == "POST" {
		f.Content = []byte(req.FormValue("body"))
		err := f.Save()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	contentType := req.Header.Get("Content-type")

	if contentType == "application/json" {
		res, _ := json.Marshal(f)
		w.Header().Set("Content-Type", "application/json")
		w.Write(res)

	} else {
		err := templates.RenderTemplate(w, "file.html", f)
		if err != nil {
			fmt.Println(err)
		}
	}

}

func (f *Textfile) Image(head []byte) error {

	file, err := os.Open(f.Abs)
	if err != nil {
		return err
	}
	defer file.Close()

	if filetype.IsImage(head) {
		fmt.Println("File is an image")
		image, _, err := image.DecodeConfig(file)
		if err != nil {
			fmt.Print(f.Path, image.ColorModel, image.Height, image.Width, err)
		}
		fmt.Print(f.Path, image.ColorModel, image.Height, image.Width, err)

		buff := make([]byte, 512) // docs tell that it take only first 512 bytes into consideration
		if _, err = file.Read(buff); err != nil {
			fmt.Println(err) // do something with that error
		}

	} else {
		fmt.Println("Not an image")
	}

	//fmt.Println(http.DetectContentType(buff)) // do something based on your detection.

	return nil
}

func (f *Textfile) Load() error {

	body, err := ioutil.ReadFile(f.Abs)
	if err != nil {
		return err
	}
	f.Title = strings.TrimSuffix(f.Name, filepath.Ext(f.Name))
	f.Content = body

	return nil
}

func (f *Textfile) Save() error {
	return ioutil.WriteFile(f.Abs, f.Content, 0600)
}
