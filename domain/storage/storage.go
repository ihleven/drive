package storage

import (
	"drive/domain"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type FileSystemStorage struct {
	Root, Prefix string
	Owner        *domain.User
	Group        *domain.Group
}

var storages = map[string]*FileSystemStorage{
	"home":   &FileSystemStorage{Root: "/Users/mi"},
	"public": &FileSystemStorage{Root: "C:/Users/ih/projects"},
}

func Register(root, prefix string) {
	storages[prefix] = &FileSystemStorage{Root: root}
}

func Get(name string) *FileSystemStorage {
	return storages[name]
}

func (st *FileSystemStorage) PermOpen(name string, flag int, uid, gid uint32) (domain.Handle, error) {
	fh, err := st.OpenFile(name, 0, 0)
	if err != nil {
		fmt.Println("PermOpen", err)
		return nil, err
	}
	//perm, err := st.GetPerm(name, uid, gid)
	//fmt.Println(err, perm)
	fmt.Println("PermOpen", name, fh)
	return fh, nil
}

func (st *FileSystemStorage) Open(name string) (domain.Handle, error) {

	return st.OpenFile(name, 0, 0)
}

func (st *FileSystemStorage) OpenFile(name string, flag int, mode os.FileMode) (domain.Handle, error) {

	path := filepath.Join(st.Root, name) // filepath.Clean(path string)

	f, err := os.OpenFile(path, flag, mode)
	if err != nil {
		return nil, err
	}

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	//var stat syscall.Stat_t
	//if err := syscall.Stat(path, &stat); err != nil {
	//	return nil, err
	//}

	return &FileHandle{
		FileInfo: info,
		File:     f,
		//Stat:     &stat,
		storage: st,
	}, nil
}

func (st *FileSystemStorage) ReadDir(path string) ([]domain.Handle, error) {

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

	entries := make([]domain.Handle, len(list))
	for index, info := range list {
		entries[index] = NewHandle(info)
	}
	return entries, nil
}

func (st *FileSystemStorage) ReadDir2(path string) ([]FileHandle, error) { // => storage

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
		//stat, _ := entry.Sys().(*syscall.Stat_t) // _ ist ok und kein error
		//fmt.Println(stat)
	}
	return handles, nil
}
