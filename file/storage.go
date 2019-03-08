package file

import (
	"drive/domain/usecase"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

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
