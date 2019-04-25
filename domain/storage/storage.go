package storage

import (
	"drive/domain"
	"drive/domain/usecase"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"drive/errors"
)

var storages = map[string]domain.Storage{
	"home":   &FileSystemStorage{Root: "/Users/mi/tmp", Prefix: "/home", Group: usecase.GetGroupByID(20)},
	"public": &FileSystemStorage{Root: "/Users/mi/Downloads", Prefix: "/public", PermissionMode: 0444},
}

func Register(root, prefix string) {
	storages[prefix] = &FileSystemStorage{Root: root}
}

func Get(name string) domain.Storage {
	return storages[name]
}

type FileSystemStorage struct {
	Root, Prefix   string
	Owner          *domain.User  // alle Dateien gehören automatisch diesem User ( => homes )
	Group          *domain.Group // jedes File des Storage bekommt automatisch diese Gruppe ( z.B. brunhilde )
	PermissionMode os.FileMode   // wenn gesetzt erhält jedes File dies Permission =< wird nicht mehr auf fs gelesen
}

func (st *FileSystemStorage) TrimPath(path string) string {
	if p := strings.TrimPrefix(path, st.Prefix); p != "" {
		return p
	}
	return "/"
}

func (st *FileSystemStorage) Location(path string) string {
	trimmedPath := strings.TrimPrefix(path, st.Prefix)
	//if trimmedPath == "" {
	//	trimmedPath = "/"
	//}
	return filepath.Join(st.Root, trimmedPath)
}

func (st *FileSystemStorage) GetHandle(name string) (domain.Handle, error) {

	location := st.Location(name)

	info, err := os.Stat(location)
	if err != nil {
		return nil, errors.Wrap(err, "os.Stat failed for %s", name)
	}

	handle := &FileHandle{
		FileInfo: info,
		mode:     info.Mode(),
		storage:  st,
		location: location,
	}

	if st.PermissionMode != 0 {
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

func (st *FileSystemStorage) ReadDir(path string) ([]domain.Handle, error) {

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

	entries := make([]domain.Handle, len(list))
	for index, info := range list {

		entries[index] = NewFileHandle(info, st, filepath.Join(location, info.Name()))
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

	file, err := os.Create(location)
	defer file.Close()
	if err != nil {
		return errors.Wrap(err, "Failed to create file: %s", location)
	}
	_, err = io.Copy(file, src)
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
