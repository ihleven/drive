package storage

import (
	"os"
	"path/filepath"
)

type Storage interface {
	Open(name string) (*FileHandle, error)
	//ReadDir(name string) ([]os.FileInfo, error)
}

var DefaultStorage *FileSystemStorage

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

func (st *FileSystemStorage) Open(name string) (*FileHandle, error) {

	//fh := FileHandle{storage: st}
	path := filepath.Join(st.Root, name) // filepath.Clean(path string)

	fd, err := os.OpenFile(path, 0, 0) // last arg: permissions when file is created
	if err != nil {
		return nil, err
	}
	// defer fd.Close()

	info, err := fd.Stat()
	if err != nil {
		return nil, err
	}

	return &FileHandle{
		FileInfo: info,
		File:     fd,
		//Stat:     &stat,
		storage: st,
	}, nil
}

func (st *FileSystemStorage) OpenFile(name string, flag int, perm os.FileMode) (*FileHandle, error) {

	path := filepath.Join(st.Root, name) // filepath.Clean(path string)

	f, err := os.OpenFile(path, flag, perm)
	if err != nil {
		return nil, err
	}

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	return &FileHandle{
		FileInfo: info,
		File:     f,
		//Stat:     &stat,
		storage: st,
	}, nil
}
