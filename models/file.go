package models

import (
	"drive/templates"
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
)

type File struct {
	*FilePath
	Directory *Directory
	MIME      types.MIME
	Title     string
	Content   []byte
}

func NewFile(fp *FilePath) *File {

	f := new(File)
	f.FilePath = fp
	f.Scan()
	return f
}

func (f *File) Render(w http.ResponseWriter, req *http.Request) {

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

func (f *File) Scan() error {
	fmt.Println("-------------------------------------------------------", f.FilePath)
	file, err := os.Open(f.FilePath.Abs)
	if err != nil {
		return err
	}
	defer file.Close()

	info, _ := file.Stat()

	ext := path.Ext(info.Name())
	if ext != "" {

		mime := mime.TypeByExtension(ext)
		fmt.Println("---------------", ext, info.Name(), mime)
		f.MIME = filetype.GetType(ext[1:]).MIME //   types.Get(ext).MIME
	}

	fmt.Printf(" - scanning '%s' - %s - %s - %s\n", f.Path, ext, f.MIME.Value)

	// We only have to pass the file header = first 261 bytes
	head := make([]byte, 261)
	file.Read(head)
	if t, e := filetype.Match(head); e == nil {
		f.MIME = t.MIME
	}

	if filetype.IsImage(head) {
		fmt.Println("File is an image")
		image, _, err := image.DecodeConfig(file)
		if err != nil {
			fmt.Print(f.Path, image.ColorModel, image.Height, image.Width, err)
		}
		fmt.Print(f.Path, image.ColorModel, image.Height, image.Width, err)

		fmt.Printf("   -> %s %s\n", info.Name())
		//d.Files = append(d.Files, val.Name())

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

func (f *File) Load() error {

	body, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return err
	}
	f.Title = strings.TrimSuffix(f.Filename, filepath.Ext(f.Filename))
	f.Content = body

	return nil
}

func (f *File) Save() error {
	return ioutil.WriteFile(f.Path, f.Content, 0600)
}
