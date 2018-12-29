package fs

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
)

type FileSystemStorage struct {
	Location string // /home/ihle
	BaseUrl  string
	//file_permissions_mode
	//directory_permissions_mode
}

func (f *FileSystemStorage) Exists(pathname string) (bool, *File) {

	fullpath := path.Join(f.Location, pathname)

	fileInfo, err := os.Stat(fullpath)
	if err != nil {
		//if os.IsNotExist(err) {
		//	return false, nil
		//}
		return false, nil
	}
	file := NewFile(&fileInfo)
	file.Path = pathname
	file.location = fullpath
	return true, file
}
func (f *FileSystemStorage) GetDirectory(pathname string) (*Directory, error) {

	location := path.Join(f.Location, pathname)

	fileInfo, err := os.Stat(location)
	if err != nil {
		return nil, err
	}
	file := NewFile(&fileInfo)
	file.Path = pathname
	file.location = location
	dir, err := NewDirectory(file)
	dir.List()
	return dir, err
}

type FileHandler interface {
	Handle(http.ResponseWriter, *http.Request)
}

func (f *FileSystemStorage) GetSpecific(pathname string) (FileHandler, error) {
	fullpath := path.Join(f.Location, pathname)

	fileInfo, error := os.Stat(fullpath)
	if error != nil {
		if os.IsNotExist(error) {
			return nil, nil
		}
		return nil, error
	}
	var fd FileHandler
	File := NewFile(&fileInfo)
	File.Path = pathname
	File.location = fullpath
	switch {
	case fileInfo.IsDir():
		File.Type = "D"
		fd, error = NewDirectory(File)

	case fileInfo.Mode().IsRegular():
		File.GuessMIME()
		File.Type = "F"
		//, error = NewFile(BaseFileInfo)
		return File, nil
	default:
		fmt.Println("It's after noon")
	}
	return fd, errors.New("Does not exist")
}

func (f *File) Parents() []File {
	var pa string
	elements := strings.Split(f.Path[1:], "/")
	list := make([]File, len(elements))
	for index, element := range elements {
		pa = fmt.Sprintf("%s/%s", pa, element)
		list[index] = File{Name: element, Path: pa}
		//fmt.Println(index, list[index].Name, list[index].Path)
	}
	return list
}

type Path struct {
	Root *string // /home/ihle
	Path string  // /alben/urlaube/Wien2013/IMG_123.jpg
	Abs  string  // /home/ihle/alben/urlaube/Wien2013/IMG_123.jpg
	//Dirname  string // /alben/urlaube/Wien2013/
	//Filename string // IMG_123.jpg
	Name     string
	Size     int64
	Mode     os.FileMode
	ModTime  time.Time
	MIME     types.MIME
	Children []Path
	Parents  []Path
}

func (p *Path) CalcParents() {
	var pa string
	elements := strings.Split(p.Path[1:], "/")
	list := make([]Path, len(elements))
	for index, element := range elements {
		pa = fmt.Sprintf("%s/%s", pa, element)
		list[index] = Path{Name: element, Path: pa}
		fmt.Println(index, list[index].Name, list[index].Path)
	}
	p.Parents = list
}

func (f *Path) MatchMIMEType() error {

	file, err := os.Open(f.Abs)
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
func (f *Path) DetectContentType() error {

	fd, err := os.Open(f.Abs)
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
func (p *Path) Type() string {
	if p.IsDir() {
		return "DIR"
	} else if p.Mode.IsRegular() {
		return "FILE"
	} else {
		return ""
	}
}

func (p *Path) IsDir() bool     { return p.Mode.IsDir() }
func (p *Path) IsRegular() bool { return p.Mode.IsRegular() }
