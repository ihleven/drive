package storage

import (
	"fmt"
	_ "image/jpeg"
	"mime"
	"net/http"
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
	mime.AddExtensionType(".md", "text/markdown")

}

func (f *FileHandle) GuessMIME() (m types.MIME) {

	if ext := path.Ext(f.Name()); ext != "" {
		mimestr := mime.TypeByExtension(ext) // mime package
		if mimestr != "" {
			m = types.NewMIME(mimestr)

		} else {
			m = filetype.GetType(ext[1:]).MIME // h2non/filetype => types.Get(ext).MIME
			//if f.MIME.Value != mimetype {
			//	fmt.Printf("MIME (%s): %s != %s \n", f.Name, mimetype, f.MIME)
			//}

		}

	}

	if m.Value == "" {
		m, _ = f.h2nonMatchMIME261()
	}
	if strings.HasSuffix(m.Subtype, "charset=utf-8") {
		m.Subtype = m.Subtype[:len(m.Subtype)-15]
		m.Value = m.Value[:len(m.Value)-15]
	}
	return
}

func (f *FileHandle) h2nonMatchMIME261() (types.MIME, error) {

	f.File.Seek(0, 0)
	// https://github.com/h2non/filetype
	// We only have to pass the file header = first 261 bytes
	head := make([]byte, 261)
	f.File.Read(head)
	t, e := filetype.Match(head)
	if e != nil {
		return t.MIME, e
	}
	//fmt.Printf(" * MIME2 '%s' => %s, %s, %s\n", f.Name, f.MIME.Value, f.MIME.Type, f.MIME.Subtype)

	return t.MIME, nil
}

func (f *FileHandle) DetectContentType() error {

	f.File.Seek(0, 0)
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := f.File.Read(buffer)
	if err != nil {
		return err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)
	m := types.NewMIME(contentType)
	fmt.Printf(" * http '%s' => %s\n", f.Name(), m.Type)

	return nil
}
