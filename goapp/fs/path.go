package fs

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
)

type FileSystemStorage struct {
	Location string // /home/ihle
	BaseUrl  string
	homes    map[string]string
	//file_permissions_mode
	//directory_permissions_mode
}

func (s *FileSystemStorage) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	file, err := s.Open(path.Clean(r.URL.Path))
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), 500)
	}

	fh, err := file.Specific()
	fh.ServeHTTP(w, r)
}
func (s *FileSystemStorage) Open(name string) (*File, error) {
	fullpath := path.Join(s.Location, name)

	fileInfo, err := os.Stat(fullpath)

	if err != nil {

		return nil, err
	}
	file := NewFile(name, fullpath, &fileInfo)

	return file, nil
}

func (f *FileSystemStorage) Exists(pathname string) (bool, *File) {

	file, err := f.Open(pathname)
	return err == nil, file
}

func (f *FileSystemStorage) OpenDir(pathname string) (*Directory, error) {

	file, err := f.Open(pathname)
	if err != nil {
		return nil, err
	}
	dir, err := NewDirectory(file)
	//dir.List()
	return dir, err
}

func (f *File) MatchMIMEType() error {

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
