package file

import (
	"drive/domain/usecase"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"syscall"
)

var storages = map[string]*FileSystemStorage{
	"home":   &FileSystemStorage{Root: "/Users/mi"},
	"public": &FileSystemStorage{Root: "/Users/mi/tmp/14"},
}
var DefaultStorage *FileSystemStorage
var PublicStorage = storages["public"]

type FileSystemStorage struct {
	Root string
	//BaseUrl  string
	//homes    map[string]string
	//file_permissions_mode
	//directory_permissions_mode
}

func SetDefaultStorage(root string) {
	DefaultStorage = &FileSystemStorage{Root: root}

}

func (st *FileSystemStorage) PermOpen(name, mode string, uid, gid uint32) (usecase.FileHandle, error) {

	fh, err := st.OpenFile(name, 0, 0)
	perm, err := st.GetPerm(name, uid, gid)
	fmt.Println(err, perm)
	return fh, nil
}

func (st *FileSystemStorage) Open(name string) (usecase.FileHandle, error) {

	return st.OpenFile(name, 0, 0)
}

func (st *FileSystemStorage) OpenFile(name string, flag int, perm os.FileMode) (usecase.FileHandle, error) {

	path := filepath.Join(st.Root, name) // filepath.Clean(path string)

	f, err := os.OpenFile(path, flag, perm)
	if err != nil {
		return nil, err
	}

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	var stat syscall.Stat_t
	if err := syscall.Stat(path, &stat); err != nil {
		return nil, err
	}

	return &FileHandle{
		FileInfo: info,
		File:     f,
		//Stat:     &stat,
		storage: *st,
	}, nil
}

func (st *FileSystemStorage) ReadDir(path string) ([]FileHandle, error) { // => storage

	f, err := os.Open(filepath.Join(st.Root, path))
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	//return list, nil

	var handles = make([]FileHandle, len(list))
	for index, entry := range list {
		if entry.Name()[0] == '.' {
			// ignore all files starting with '.'
			continue
		}
		handles[index] = FileHandle{FileInfo: entry}
		stat, _ := entry.Sys().(*syscall.Stat_t) // _ ist ok und kein error
		fmt.Println(stat)
	}
	return handles, nil
}
