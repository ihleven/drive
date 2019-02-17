package file

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"syscall"
	"time"
)

type Storage interface {
	Open(name string) (os.File, error)
	GetFile(name string) (*File, error)
}

type FileSystemStorage struct {
	Root string
}

var storage = &FileSystemStorage{Root: "/Users/mi/go"}

type Info struct {
	os.FileInfo
	//File    *os.File
	Descriptor *os.File
	sysStat    *syscall.Stat_t
}

func (s *FileSystemStorage) Open(path string) (*os.File, error) {

	fullpath := filepath.Join(s.Root, path)

	file, err := os.Open(fullpath)
	if err != nil {
		return nil, err
	}

	return file, nil
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
func (s *FileSystemStorage) Stat(path string, uid, gid uint32) (*Info, error) {

	fullpath := filepath.Join(s.Root, path)

	fileInfo, err := os.Stat(fullpath)
	if err != nil {
		return nil, err
	}

	file := &Info{
		FileInfo: fileInfo,
		sysStat:  fileInfo.Sys().(*syscall.Stat_t),
	}
	//file := NewFile(name, fullpath, &fileInfo)
	return file, nil
}

func (f *Info) Permissions(uid, gid uint32) (r, w, x bool) {

	f.sysStat = f.Sys().(*syscall.Stat_t)
	//fmt.Printf("sysstat %s\n", f.sysStat)

	if f.sysStat == nil {
		return false, false, false
	}

	//uid, err := strconv.Atoi(usr.Uid)
	switch {
	case f.sysStat.Uid == uid:
		return f.UserPermissions()
	case f.sysStat.Gid == gid:
		return f.GroupPermissions()
	default:
		return f.OthersPermissions()
	}
}
func (h *Info) UserPermissions() (r, w, x bool) {
	rwx := fmt.Sprintf("%s", h.FileInfo.Mode())
	fmt.Println("UserPermission:", rwx, rwx[1] == 'r', rwx[2] == 'w', rwx[3] == 'x')
	return rwx[1] == 'r', rwx[2] == 'w', rwx[3] == 'x'
}
func (h *Info) GroupPermissions() (r, w, x bool) {
	rwx := fmt.Sprintf("%s", h.FileInfo.Mode())
	fmt.Println("GroupPermission:", rwx, rwx[4] == 'r', rwx[5] == 'w', rwx[6] == 'x')
	return rwx[4] == 'r', rwx[5] == 'w', rwx[6] == 'x'
}
func (h *Info) OthersPermissions() (r, w, x bool) {
	rwx := fmt.Sprintf("%s", h.FileInfo.Mode())
	fmt.Println(rwx, rwx[7] == 'r', rwx[8] == 'w', rwx[9] == 'x')
	return rwx[7] == 'r', rwx[8] == 'w', rwx[9] == 'x'
}

// func (f *File) GetPerm(fi *os.FileInfo) error {

// 	stat := (*fi).Sys().(*syscall.Stat_t)
// 	uid := fmt.Sprintf("%d", stat.Uid)
// 	gid := fmt.Sprintf("%d", stat.Gid)
// 	f.Group, _ = user.LookupGroupId(gid)
// 	f.User, _ = user.LookupId(uid)
// 	return nil
// }
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
