package file

import (
	"fmt"
	_ "image/jpeg"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
)

func init() {
	mime.AddExtensionType(".py", "text/python")
	mime.AddExtensionType(".go", "text/golang")
	mime.AddExtensionType(".json", "text/json")
	mime.AddExtensionType(".js", "text/javascript")
	mime.AddExtensionType(".ts", "text/typescript")
	mime.AddExtensionType(".dia", "text/diary")

}

func (f *File) Specific() error {

	if f.Mode.IsRegular() {
		f.Type = "F"
		f.GuessMIME()
		switch f.MIME.Type {
		case "image":
			//fh, err = f.AsImage()
			//image.ServeHTTP(w, r)

		case "text":
			switch f.MIME.Subtype {
			case "diary; charset=utf-8":
				//fh, err = NewDiary(f)
				//fh, err = f.AsTextfile()
			default:
				//fh, err = f.AsTextfile()
				//fmt.Println("TESTFILE", fh, f.MIME.Subtype)
			}

		default:
			fmt.Println("No specific file type found")
			//fh, err = f, nil
		}

	}
	return nil
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

	}

	if f.MIME.Value == "" {
		f.h2nonMatchMIME261()
	}
	if strings.HasSuffix(f.MIME.Subtype, "charset=utf-8") {
		f.MIME.Subtype = f.MIME.Subtype[:len(f.MIME.Subtype)-15]
		f.MIME.Value = f.MIME.Value[:len(f.MIME.Value)-15]
	}
	fmt.Println(f.MIME)
	switch f.MIME.Type {
	case "image":
		f.Type = "FI"
	case "text":
		f.Type = "FT"
	default:
		f.Type = "F"
	}

}
func (f *File) h2nonMatchMIME261() error {

	file, err := os.Open(f.location)
	if err != nil {
		return err
	}
	defer file.Close()

	// https://github.com/h2non/filetype
	// We only have to pass the file header = first 261 bytes
	head := make([]byte, 261)
	file.Read(head)
	if t, e := filetype.Match(head); e == nil {
		f.MIME = t.MIME
	}
	fmt.Printf(" * MIME2 '%s' => %s, %s, %s\n", f.Name, f.MIME.Value, f.MIME.Type, f.MIME.Subtype)

	return nil
}

func (f *File) DetectContentType() error {

	fd, err := os.Open(f.location)
	if err != nil {
		return err
	}
	defer fd.Close()

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err = fd.Read(buffer)
	if err != nil {
		return err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)
	f.MIME = types.NewMIME(contentType)
	fmt.Printf(" * http '%s' => %s\n", f.Name, f.MIME.Type)

	return nil
}
