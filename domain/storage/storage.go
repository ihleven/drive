package storage

import (
	"drive/domain"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var storages = map[string]*FileSystemStorage{
	"home":   &FileSystemStorage{Root: "/Users/mi/tmp", Prefix: "/home"},
	"public": &FileSystemStorage{Root: "/Users/mi/Downloads", Prefix: "/public", PermissionMode: 0444},
}

func Register(root, prefix string) {
	storages[prefix] = &FileSystemStorage{Root: root}
}

func Get(name string) *FileSystemStorage {
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

func (st *FileSystemStorage) GetHandle(name string) (domain.Handle, error) {

	location := filepath.Join(st.Root, st.TrimPath(name))

	info, err := os.Stat(location)
	if err != nil {
		return nil, err
	}
	handle := &FileHandle{
		FileInfo: info,
		storage:  st,
		location: location,
		mode:     info.Mode(),
	}

	if st.PermissionMode != 0 {
		//fmt.Printf("Mode %b %b %b \n\n", handle.mode, handle.mode&0xfffffe00, st.PermissionMode&0x1ff)
		handle.mode = (handle.mode & 0xfffffe00) | (st.PermissionMode & 0x1ff)
		//fmt.Printf("Mode %b  \n\n", handle.mode)
	}
	return handle, nil
}

func (st *FileSystemStorage) OpenFD(name string) (*os.File, error) {
	path := filepath.Join(st.Root, st.TrimPath(name))
	fmt.Println("OpenFD", name, st.TrimPath(name))
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return fd, nil
}

func (st *FileSystemStorage) ReadDir(name string) ([]domain.Handle, error) {

	location := filepath.Join(st.Root, st.TrimPath(name))

	fd, err := os.Open(location)
	defer fd.Close()
	if err != nil {
		return nil, err
	}

	list, err := fd.Readdir(-1)
	if err != nil {
		return nil, err
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })

	entries := make([]domain.Handle, 0)
	for _, info := range list {
		handle := &FileHandle{FileInfo: info, storage: st, mode: info.Mode(), location: filepath.Join(location, info.Name())}

		if st.PermissionMode != 0 {
			//fmt.Printf("Mode %b %b %b \n\n", handle.mode, handle.mode&0xfffffe00, st.PermissionMode&0x1ff)
			handle.mode = (handle.mode & 0xfffffe00) | (st.PermissionMode & 0x1ff)
			//fmt.Printf("Mode %b  \n\n", handle.mode)
		}
		entries = append(entries, handle)
	}
	return entries, nil
}
