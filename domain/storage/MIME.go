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

	if f.IsDir() {
		m = types.NewMIME("directory/")
		return
	}

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
		// m, _ = f.h2nonMatchMIME261()
	}
	if strings.HasSuffix(m.Subtype, "charset=utf-8") {
		m.Subtype = m.Subtype[:len(m.Subtype)-15]
		m.Value = m.Value[:len(m.Value)-15]
	}
	return
}

func (fh *FileHandle) h2nonMatchMIME261() (types.MIME, error) {
	fd := fh.Descriptor()
	defer fd.Close()
	fd.Seek(0, 0)
	// https://github.com/h2non/filetype
	// We only have to pass the file header = first 261 bytes
	head := make([]byte, 261)
	fd.Read(head)
	typ, err := filetype.Match(head)
	if err != nil {
		return typ.MIME, err
	}
	return typ.MIME, nil
}

func (fh *FileHandle) DetectContentType() error {
	fd := fh.Descriptor()
	defer fd.Close()
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := fd.Read(buffer)
	if err != nil {
		return err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)
	m := types.NewMIME(contentType)
	fmt.Printf(" * http '%s' => %s\n", fh.Name(), m.Type)

	return nil
}
