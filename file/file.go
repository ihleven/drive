package file

import (
	"drive/auth"
	"errors"
	"fmt"
	"os"
	"os/user"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/h2non/filetype/types"
)

type File struct {
	// location   string
	// Path       string `json:"path"`
	//Name string `json:"name"`
	// Size       int64  `json:"size"`
	//Mode os.FileMode
	//Modtime time.Time `json:"mtime"`
	// AccessTime time.Time
	// ChangeTime time.Time
	// MIME       types.MIME
	// Type       string `json:"type"`
	// User       *user.User
	// Group      *user.Group
	// Uid        int
	// Gid        int
	*Info
	location string
	Path     string `json:"path"`
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Mode     os.FileMode
	ModTime  time.Time `json:"mtime"`
	MIME     types.MIME
	Type     string `json:"type"`
	User     *user.User
	Group    *user.Group
}

func NewFile(path string, usr *auth.User) (*File, error) {

	info, err := storage.GetFile(path)
	if err != nil {
		return nil, err
	}

	f := &File{
		Info: info,
		Path: path,
		Mode: info.Mode(),
		Name: info.Name(),
		Size: info.Size(),
	}

	return f, nil
}

func (f *File) GetContent() ([]byte, error) { //offset, limit int) (e error) {

	var buffer = make([]byte, f.Size)
	len, err := f.Descriptor.Read(buffer)
	if err != nil {
		return nil, err
	}

	if int64(len) != f.Size {
		return buffer, errors.New(fmt.Sprintf("read only %d of %d bytes", len, f.Size))
	}
	return buffer, nil
}

func (f *File) GetTextContent() (string, error) { //offset, limit int) (e error) {

	var body = make([]byte, f.Size)
	len, err := f.Descriptor.Read(body)
	if err != nil {
		return "", err
	}
	if int64(len) != f.Size {
		return "", errors.New(fmt.Sprintf("read only %d of %d bytes", len, f.Size))
	}
	if utf8.Valid(body) {
		return string(body), nil
	} else {
		return "", errors.New("Invalid UTF-8")
	}

}

func (f *File) SetContent(content string) error { //offset, limit int) (e error) {

	//body, e := ioutil.ReadFile(f.location)
	//if utf8.Valid(body) {
	//	f.Content = body
	//} else {
	//	e = errors.New("Invalid UTF-8")
	//}
	return nil
}

func (f *File) Breadcrumbs() []map[string]string {

	////////////////////

	var _ = func(path string) (dir, file string) {
		i := strings.LastIndex(path, "/")
		return path[:i+1], path[i+1:]
	}

	elements := strings.Split(strings.Trim(f.Path[1:], "/"), "/")
	breadcrumbs, currentPath := make([]map[string]string, len(elements)), ""
	for index, element := range elements {
		currentPath = currentPath + "/" + element
		breadcrumbs[index] = map[string]string{"name": element, "path": currentPath} // "/" + strings.Join(elements[:index+1], "/")}
	}
	breadcrumbs[len(elements)-1]["path"] = ""
	return breadcrumbs
}

func (f *File) Parents() []struct{ Name, Path string } {

	var path string
	elements := strings.Split(f.Path[1:], "/")
	list := make([]struct{ Name, Path string }, len(elements))
	for index, element := range elements {
		path = fmt.Sprintf("%s/%s", path, element)
		list[index] = struct{ Name, Path string }{Name: element, Path: path}
	}
	return list
}
func (f *File) FormattedMTime() string {

	return f.ModTime.Format(time.RFC822Z)
}
