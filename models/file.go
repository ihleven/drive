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

type FilePath struct {
	Info     os.FileInfo
	Abs      string // /home/ihle/alben/urlaube/Wien2013/IMG_123.jpg
	Path     string // /alben/urlaube/Wien2013/IMG_123.jpg
	Route    string // /alben/urlaube/Wien2013/IMG_123.jpg
	Dirname  string // /alben/urlaube/Wien2013/
	Filename string // IMG_123.jpg

	Root string // /home/ihle
}

func GetBaseFile(url string) *BaseFile {

	fp := path.Join(Root_folder, url)

	info, error := os.Stat(fp)
	if error != nil {
		return nil
	}
	return &BaseFile{fp, info}
}

type File struct {
	FilePath
	Directory *Directory
	MIME      types.MIME
	Title     string
	Body      []byte
}

func NewFile(fileInfo os.FileInfo, url string) (*File, error) {
	abspath := path.Join(Root_folder, url)
	dirname, basename := filepath.Split(url)

	return &File{Path: abspath, Name: basename, Directory: dirname}, nil
}

func (f *File) Render(w http.ResponseWriter, req *http.Request) {

	contentType := req.Header.Get("Content-type")

	if contentType != "application/json" {
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

	file, err := os.Open(f.Path)
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
	f.Title = strings.TrimSuffix(f.Name, filepath.Ext(f.Name))
	f.Body = body

	return nil
}

func (f *File) Save() error {
	return ioutil.WriteFile(f.Path, f.Body, 0600)
}
