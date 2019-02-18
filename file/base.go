package file

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"syscall"
	"time"
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
}

var storage = &FileSystemStorage{Root: "/Users/mi/go"}

type FileSystemStorage struct {
	Root string
	//BaseUrl  string
	//homes    map[string]string
	//file_permissions_mode
	//directory_permissions_mode
}

func (s *FileSystemStorage) Open(path string, flag int, uid, gid uint32) (*Info, error) {

	fullpath := filepath.Join(s.Root, path)

	fd, err := os.OpenFile(fullpath, flag, 0)
	if err != nil {
		return nil, err
	}

	info, err := fd.Stat()
	if err != nil {
		return nil, err
	}

	//stat = info.Sys().(*syscall.Stat_t)
	var stat syscall.Stat_t
	if err := syscall.Stat(fullpath, &stat); err != nil {
		return nil, err
	}

	var mode = stat.Mode
	var shift uint32

	switch {
	case stat.Uid == uid:
		shift = 6
	case stat.Gid == gid:
		shift = 3
	default:
		shift = 0
	}
	var r, w, x = mode & (OS_READ << shift), mode & (OS_WRITE << shift), mode & (OS_EX << shift)

	if (flag == os.O_RDONLY && r == 0) || (flag == os.O_WRONLY && w == 0) {
		return nil, os.ErrPermission
	}
	fmt.Printf("RWX: %d     %t %t %t \n", shift, r != 0, w != 0, x != 0)

	return &Info{
		FileInfo:   info,
		Descriptor: fd,
		Stat:       &stat,
	}, nil
}
func (s *FileSystemStorage) GetFile(path string) (*Info, error) {

	fullpath := filepath.Join(s.Root, path)

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

type Info struct {
	os.FileInfo
	Descriptor *os.File
	Stat       *syscall.Stat_t
}

func (f *Info) GetPermissions(uid, gid uint32) (r, w, x bool) {

	fmt.Println(uid, gid)
	// f.Stat = f.Sys().(*syscall.Stat_t)

	if f.Stat == nil {
		return false, false, false
	}

	switch {
	case f.Stat.Uid == uid:
		return f.UserPermissions()
	case f.Stat.Gid == gid:
		return f.GroupPermissions()
	default:
		return f.OthersPermissions()
	}
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

func statAtime(st *syscall.Stat_t) time.Time {
	return time.Unix(st.Atimespec.Unix())
}
func statCtime(st *syscall.Stat_t) time.Time {
	return time.Unix(st.Ctimespec.Unix())
}

// ReadDir reads the directory named by dirname and returns
// a list of directory entries sorted by filename.
func (i *Info) ReadDir() ([]os.FileInfo, error) {

	list, err := i.Descriptor.Readdir(-1)
	if err != nil {
		return nil, err
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	return list, nil
}

func (i *Info) Close() error {
	err := i.Descriptor.Close()
	return err
}
