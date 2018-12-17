package fs

import (
	"drive/templates"
	"encoding/json"
	"fmt"
	"net/http"
)

type Filer interface {
	Render(w http.ResponseWriter, r *http.Request)
	//ContentType() string
}

func NewFiler(root, pathname string) (Filer, error) {

	basefile, err := NewPath(root, pathname)
	if basefile == nil {
		return nil, err
	}
	if basefile.IsDir() {
		return nil, nil // NewDirectory(basefile)
	} else {
		return NewFile(basefile), nil
	}

}

type Directory struct {
	Path
	Parent    string
	Folders   []string
	Files     []*Path
	IndexFile string
}

func NewDirectory(fp *Path) (*Directory, error) {
	d := &Directory{Path: *fp}
	//d.scan()
	return d, nil
}

func (d *Directory) Render(w http.ResponseWriter, req *http.Request) {

	if contentType := req.Header.Get("Content-type"); contentType == "application/json" {
		json.NewEncoder(w).Encode(d)
	} else {
		fmt.Println("asdfasdfasdf")
		templates.RenderTemplate(w, "directory.html", d)
	}
}
