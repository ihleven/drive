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

func (st *FileSystemStorage) trimPath(path string) string {
	if p := strings.TrimPrefix(path, st.BaseURL); p != "" {
		return p
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

	trimmedPath := strings.TrimPrefix(path, st.BaseURL)
	//if trimmedPath == "" {
	//	trimmedPath = "/"
	//}
	if trimmedPath != path {
		fmt.Printf("trimmed (%s): %s => %s, location: %s \n", st.BaseURL, path, trimmedPath, filepath.Join(st.Root, trimmedPath))
	} else {
		//fmt.Printf("not trimmed: %s => %s\n", path, filepath.Join(st.Root, trimmedPath))
	}
	return filepath.Join(st.Root, trimmedPath)
}

func (st *FileSystemStorage) GetHandle(url string) (drive.Handle, error) {

	path := strings.TrimPrefix(url, st.BaseURL)
	location := filepath.Join(st.Root, path)

	info, err := os.Stat(location)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.Augment(err, errors.NotFound, "os.Stat failed for %s (location: %s)", url, location)
		}
		return nil, errors.Wrap(err, "os.Stat failed for %s (location: %s)", url, location)
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

func (st *FileSystemStorage) Save(path string, src io.Reader) error { // content []byte

	location := st.Location(path)
	// detect if file exists
	var _, err = os.Stat(location)
	if err != nil && !os.IsNotExist(err) {
		return errors.Wrap(err, "os.Stat error for: %s", path)
	}
	if err == nil {
		return errors.Errorf("File cannot be created because it already exists %s", location)
	}

	dest, err := os.Create(location)
	defer dest.Close()
	if err != nil {
		return errors.Wrap(err, "Failed to create file: %s", location)
	}
	_, err = io.Copy(dest, src)
	if err != nil {
		return errors.Wrap(err, "Failed to copy content to file")
	}
	return nil
}

func (st *FileSystemStorage) Create(path string) error {
	fmt.Println("storage create", path)
	location := st.Location(path)
	var _, err = os.Stat(location)
	if err != nil && !os.IsNotExist(err) {
		return errors.Wrap(err, "os.Stat error for: %s", path)
	}
	if err == nil {
		return errors.Errorf("File cannot be created because it already exists %s", location)
	}

	file, err := os.Create(location)
	defer file.Close()
	if err != nil {
		return errors.Wrap(err, "Failed to create file: %s", path)
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
