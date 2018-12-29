package fs

import (
	"fmt"
	_ "image/jpeg"
	"mime"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
)

type File struct {
	location string
	Path     string `json:"path"`
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Mode     os.FileMode
	ModTime  time.Time `json:"mtime"`
	MIME     types.MIME
	Type     string `json:"type"`
}

func NewFile(fi *os.FileInfo) *File {

	return &File{Name: (*fi).Name(), Size: (*fi).Size(), Mode: (*fi).Mode(), ModTime: (*fi).ModTime()}
}

func (f *File) String() string {
	return fmt.Sprintf("%s: %s", f.Type, f.Path)
}
func (f *File) FormattedMTime() string {
	return f.ModTime.Format(time.RFC822Z)
}

func (f *File) GuessMIME() {

	if ext := path.Ext(f.Name); ext != "" {
		mime := mime.TypeByExtension(ext) // mime package
		if mime != "" {
			f.MIME = types.NewMIME(mime)

		} else {
			f.MIME = filetype.GetType(ext[1:]).MIME // h2non/filetype => types.Get(ext).MIME
			//if f.MIME.Value != mimetype {
			//	fmt.Printf("MIME (%s): %s != %s \n", f.Name, mimetype, f.MIME)
			//}

		}
		if f.MIME.Value == "" {
			file, _ := os.Open(f.location)

			// We only have to pass the file header = first 261 bytes
			head := make([]byte, 261)
			file.Read(head)
			t, _ := filetype.Get(head)
			f.MIME = t.MIME

		}

	}
	switch f.MIME.Type {
	case "image":
		f.Type = "FI"
	case "text":
		f.Type = "FT"
	default:
		f.Type = "F"
	}

}
func (f *File) Handle(w http.ResponseWriter, r *http.Request) {

}
