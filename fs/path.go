package fs

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"path"
)

func init() {
	mime.AddExtensionType(".py", "text/python")
	mime.AddExtensionType(".go", "text/golang")
	mime.AddExtensionType(".json", "text/json")
	mime.AddExtensionType(".js", "text/javascript")
	mime.AddExtensionType(".ts", "text/typescript")
	mime.AddExtensionType(".dia", "text/diary")

}

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
			fmt.Println("asdf =>", s.Location)
			http.ServeFile(w, r, "../vue/dist/index.html")
			//http.NotFound(w, r)
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

func (f *FileSystemStorage) OpenDir(pathname string) (*Directory, error) {

	file, err := f.Open(pathname)
	if err != nil {
		return nil, err
	}
	dir, err := NewDirectory(file)
	//dir.List()
	return dir, err
}
