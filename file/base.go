package file

import (
	"drive/config"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

const (
	OS_READ        = 04
	OS_WRITE       = 02
	OS_EX          = 01
	OS_USER_SHIFT  = 6
	OS_GROUP_SHIFT = 3
	OS_OTH_SHIFT   = 0

	OS_USER_R   = OS_READ << OS_USER_SHIFT
	OS_USER_W   = OS_WRITE << OS_USER_SHIFT
	OS_USER_X   = OS_EX << OS_USER_SHIFT
	OS_USER_RW  = OS_USER_R | OS_USER_W
	OS_USER_RWX = OS_USER_RW | OS_USER_X

	OS_GROUP_R   = OS_READ << OS_GROUP_SHIFT
	OS_GROUP_W   = OS_WRITE << OS_GROUP_SHIFT
	OS_GROUP_X   = OS_EX << OS_GROUP_SHIFT
	OS_GROUP_RW  = OS_GROUP_R | OS_GROUP_W
	OS_GROUP_RWX = OS_GROUP_RW | OS_GROUP_X

	OS_OTH_R   = OS_READ << OS_OTH_SHIFT
	OS_OTH_W   = OS_WRITE << OS_OTH_SHIFT
	OS_OTH_X   = OS_EX << OS_OTH_SHIFT
	OS_OTH_RW  = OS_OTH_R | OS_OTH_W
	OS_OTH_RWX = OS_OTH_RW | OS_OTH_X

	OS_ALL_R   = OS_USER_R | OS_GROUP_R | OS_OTH_R
	OS_ALL_W   = OS_USER_W | OS_GROUP_W | OS_OTH_W
	OS_ALL_X   = OS_USER_X | OS_GROUP_X | OS_OTH_X
	OS_ALL_RW  = OS_ALL_R | OS_ALL_W
	OS_ALL_RWX = OS_ALL_RW | OS_GROUP_X
)

type Storage interface {
	Open(name string) (os.File, error)
	GetFile(name string) (*File, error)
	ReadDir(name string) ([]os.FileInfo, error)
}

var storage = &FileSystemStorage{Root: &config.Root}

type FileSystemStorage struct {
	Root *string
	//BaseUrl  string
	//homes    map[string]string
	//file_permissions_mode
	//directory_permissions_mode
}

func Open(path string, flag int, uid, gid uint32) (*Info, error) {

	fullpath := filepath.Join(*storage.Root, path)
	fmt.Println(fullpath, flag)
	fd, err := os.OpenFile(fullpath, flag, 0)
	if err != nil {
		return nil, err
	}

	info, err := fd.Stat()
	if err != nil {
		return nil, err
	}

	return &Info{
		FileInfo:   info,
		Descriptor: fd,
	}, nil
}
func (s *FileSystemStorage) GetFile(path string) (*Info, error) {

	fullpath := filepath.Join(*s.Root, path)

	fd, err := os.Open(fullpath)
	if err != nil {
		return nil, err
	}

	stat, err := fd.Stat()
	if err != nil {
		return nil, err
	}
	file := &Info{FileInfo: stat, Descriptor: fd}

	return file, nil
}

func ReadDir(path string) ([]os.FileInfo, error) {

	f, err := os.Open(filepath.Join(*storage.Root, path))
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	return list, nil
}

func Mkdir(path string) {
	fullpath := filepath.Join(*storage.Root, path)
	os.MkdirAll(fullpath, os.ModePerm)

}

type Info struct {
	os.FileInfo
	Descriptor *os.File
}

func (f *Info) GetPermissions(uid, gid uint32) (r, w, x bool) {

	//fmt.Println(uid, gid)
	// f.Stat = f.Sys().(*syscall.Stat_t)
	return r, w, x
}
func (f *Info) UserPermissions() (r, w, x bool) {
	rwx := f.Mode().String()
	return rwx[1] == 'r', rwx[2] == 'w', rwx[3] == 'x'
}
func (f *Info) GroupPermissions() (r, w, x bool) {
	rwx := f.Mode().String()
	return rwx[4] == 'r', rwx[5] == 'w', rwx[6] == 'x'
}
func (f *Info) OthersPermissions() (r, w, x bool) {
	rwx := f.Mode().String()
	return rwx[7] == 'r', rwx[8] == 'w', rwx[9] == 'x'
}

// ReadDir reads the directory named by dirname and returns
// a list of directory entries sorted by filename.
// from os.File.Readdir
func (i *Info) ReadDir() ([]os.FileInfo, error) {

	list, err := i.Descriptor.Readdir(-1)
	if err != nil {
		return nil, err
	}
	//sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	return list, nil
}

func (i *Info) Close() error {
	err := i.Descriptor.Close()
	return err
}
