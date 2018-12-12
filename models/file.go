package models

import (
	"drive/templates"
	"encoding/json"
	"fmt"
	"image"
	"mime"
	"net/http"
	"os"
	"path"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
)

type File struct {
	basefile
	Path      string
	Name      string
	Directory *Directory
	MIME      types.MIME
	Title     string
	Body      []byte
}

func (f *File) Render(w http.ResponseWriter, req *http.Request) {

	contentType := req.Header.Get("Content-type")

	if contentType != "application/json" {
		res, _ := json.Marshal(f)
		w.Header().Set("Content-Type", "application/json")
		w.Write(res)

	} else {
		err := templates.Templates.ExecuteTemplate(w, "file.html", f)
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
