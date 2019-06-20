package storage

import (
	"drive/domain"
	"drive/drive"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"drive/errors"
)

var storages = map[string]drive.Storage{
	//"home": &FileSystemStorage{Root: "/Users/mi/tmp", Prefix: "/home", Group: usecase.GetGroupByID(20)},
	//"public": &FileSystemStorage{Root: "/Users/mi/Downloads", Prefix: "/public", PermissionMode: 0444},
}

func Register(name, root, baseUrl, serveUrl string, group *drive.Group, mode os.FileMode) drive.Storage {

	storages[name] = &FileSystemStorage{
		Root:           root,
		BaseURL:        baseUrl,
		ServeURL:       serveUrl,
		Group:          group,
		PermissionMode: mode,
	}
	return storages[name]
}

func Get(name string) drive.Storage {
	return storages[name]
}

type FileSystemStorage struct {
	Root           string          `json:"-"`
	BaseURL        string          `json:"baseUrl"`
	ServeURL       string          `json:"serveUrl"`
	AlbumURL       string          `json:"albumUrl"`
	Owner          *drive.User     `json:"-"` // alle Dateien gehören automatisch diesem User ( => homes )
	Group          *drive.Group    `json:"-"` // jedes File des Storage bekommt automatisch diese Gruppe ( z.B. brunhilde )
	PermissionMode os.FileMode     `json:"-"` // wenn gesetzt erhält jedes File dies Permission =< wird nicht mehr auf fs gelesen
	Account        *domain.Account `json:"-"` //
}

func (st *FileSystemStorage) CleanPath(path string) string {

	cleanedPath := strings.TrimPrefix(filepath.Clean(path), st.BaseURL)
	if cleanedPath != path {
		fmt.Printf("CleanPath(%s) => %s, BaseURL: %s \n", path, cleanedPath, st.BaseURL)
		//return trimmedPath, errors.New(errors.PathError, "Path was trimmed")
	}
	if cleanedPath != "" {
		return cleanedPath
	}
	return "/"
}

func (st *FileSystemStorage) CleanServePath(path string) string {

	cleanedPath := strings.TrimPrefix(filepath.Clean(path), st.ServeURL)
	if cleanedPath != path {
		fmt.Printf("CleanPath(%s) => %s, BaseURL: %s \n", path, cleanedPath, st.ServeURL)
		//return trimmedPath, errors.New(errors.PathError, "Path was trimmed")
	}
	if cleanedPath != "" {
		return cleanedPath
	}
	return "/"
}

func (st *FileSystemStorage) URL(path string) string {
	return filepath.Join(st.BaseURL, path)
}
func (st *FileSystemStorage) GetServeURL(path string) string {
	return filepath.Join(st.ServeURL, path)
}

func (st *FileSystemStorage) Location(path string) string {

	return filepath.Join(st.Root, path)
}

func (st *FileSystemStorage) GetHandle(path string) (drive.Handle, error) {

	info, err := os.Stat(filepath.Join(st.Root, path))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.Augment(err, errors.NotFound, "os.Stat failed for %s (location: %s)", path, filepath.Join(st.Root, path))
		}
		return nil, errors.Wrap(err, "os.Stat failed for %s (location: %s)", path, filepath.Join(st.Root, path))
	}

	handle := &FileHandle{
		FileInfo: info,
		storage:  st,
		path:     path,
		mode:     info.Mode(),
	}

	if st.PermissionMode != 0 {
		// replace 9 least significant bits from mode with storage.PermissionMode
		handle.mode = (handle.mode & 0xfffffe00) | (st.PermissionMode & 0x1ff)
	}
	return handle, nil
}

// Create creates the named file with mode 0666 (before umask), truncating it if it already exists.
// If successful, methods on the returned File can be used for I/O;
// the associated file descriptor has mode O_RDWR. If there is an error, it will be of type *PathError.
func (st *FileSystemStorage) Create(path string, overwrite bool) (*os.File, error) {

	location := st.Location(path)

	if !overwrite {
		// check if file exists
		if exists, err := st.Exists(path); err != nil && exists {
			return nil, errors.Errorf("File cannot be created because it already exists %s", location)
		}
		// file existiert nicht oder keine Berechtigung oder ...
		// jedenfalls kann os.Create aufgerufen werden

		if _, err := os.Stat(location); err == nil {
			// pfad zu location exists
			return nil, errors.Errorf("File cannot be created because it already exists %s", location)
		}
	}
	fd, err := os.Create(location)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create file %v", location)
	}
	return fd, nil
}
func (st *FileSystemStorage) Exists(path string) (bool, error) {

	if _, err := os.Stat(st.Location(path)); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		return false, errors.Wrap(err, "Schrödinger: file '%s' may or may not exist.", path)
	}
}

// os.Open:
// Open opens the named file for reading. If successful, methods on
// the returned file can be used for reading; the associated file
// descriptor has mode O_RDONLY.
// If there is an error, it will be of type *PathError.
//func Open(name string) (*File, error) {
//	return OpenFile(name, O_RDONLY, 0)
//}
func (st *FileSystemStorage) Open(path string) (*os.File, error) {

	location := st.Location(path)
	fd, err := os.Open(location)
	if err != nil {
		return nil, errors.Wrap(err, "Could not get file descriptor for %v", location)
	}
	return fd, nil
}

func (st *FileSystemStorage) ReadDir(path string) ([]drive.Handle, error) {

	location := st.Location(path)
	fd, err := os.Open(location)
	defer fd.Close()
	if err != nil {
		return nil, errors.Wrap(err, "Could not get file descriptor for %v", location)
	}

	list, err := fd.Readdir(-1)
	if err != nil {
		return nil, errors.Wrap(err, "Could not read dir %v", fd)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })

	entries := make([]drive.Handle, len(list))
	for index, info := range list {

		entries[index] = NewFileHandle(info, st, filepath.Join(path, info.Name()))
	}
	return entries, nil
}

func (st *FileSystemStorage) Save(path string, content io.Reader, overwrite bool) error { // content []byte

	name := st.Location(path)
	var file *os.File
	var err error

	if overwrite {
		// create a new file if none exists, truncate existing file when opened.
		file, err = os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	} else {
		// create a new file, file must not exist.
		file, err = os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	}
	if err != nil {
		return errors.Wrap(err, "Could not create file: %s (overwrite: %v)", path, overwrite)
	}
	defer file.Close()

	_, err = io.Copy(file, content)
	if err != nil {
		return errors.Wrap(err, "Failed to save content to file %v", path)
	}
	return nil
}

func (st *FileSystemStorage) Delete(path string) error {
	fmt.Println("storage delete", path)
	location := st.Location(path)
	var err = os.Remove(location)
	if err != nil {
		return errors.Wrap(err, "Failed to delete %s", path)
	}
	return nil
}
