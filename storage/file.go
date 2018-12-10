package storage

import (
	"fmt"
	"html/template"
	"image"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
)

type BaseFile struct {
	Path string
	Info os.FileInfo
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
	Info      os.FileInfo
	Url       string
	Path      string
	Name      string
	Directory string
	MIME      types.MIME
	Title     string
	Body      []byte
}

func NewFile(fileInfo os.FileInfo, url string) (*File, error) {
	abspath := path.Join(Root_folder, url)
	dirname, basename := filepath.Split(url)

	return &File{Path: abspath, Name: basename, Directory: dirname}, nil
}

func (f *File) Scan() error {

	file, err := os.Open(f.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	info, _ := file.Stat()

	ext := path.Ext(info.Name())
	mime := mime.TypeByExtension(ext)
	f.MIME = filetype.GetType(ext[1:]).MIME //   types.Get(ext).MIME

	fmt.Printf(" - scanning '%s' - %s - %s - %s\n", f.Path, ext, mime, f.MIME.Value)

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

		fmt.Printf("   -> %s %s\n", info.Name(), mime)
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
	println(f.Path)

	kind, e := filetype.Match(body)
	if e == nil {
		println("nil")
	}
	println("asdf", kind.Extension)
	f.MIME = kind.MIME
	//if f.MIME.Type == "text" {
	f.Title = strings.TrimSuffix(f.Name, filepath.Ext(f.Name))
	f.Body = body
	//}

	return nil
}

func (f *File) Save() error {
	return ioutil.WriteFile(f.Path, f.Body, 0600)
}

func (f *File) RenderHTML(w http.ResponseWriter, req *http.Request) {

	tpl, err := template.ParseFiles("templates/directory.html")
	if err != nil {
		http.Error(w, "500 Internal Error : Error while generating directory listing.", 500)
		fmt.Println(err)
		log.Print(err)
		return
	}

	err = tpl.Execute(w, f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (d *File) MarschalJSON(w http.ResponseWriter, req *http.Request) {

}
