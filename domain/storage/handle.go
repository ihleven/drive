package storage

import (
	"drive/domain"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

type FileHandle struct {
	os.FileInfo
	File *os.File
	//Stat    *syscall.Stat_t
	storage domain.Storage
}

func NewHandle(info os.FileInfo) domain.Handle { // https://stackoverflow.com/questions/40823315/x-does-not-implement-y-method-has-a-pointer-receiver
	h := FileHandle{FileInfo: info}
	return &h
}

func (fh *FileHandle) GetFile() *os.File {
	return fh.File
}

func (fh *FileHandle) Close() error {
	if fh.File != nil {
		err := fh.File.Close()
		fh.File = nil
		return err
	}
	return errors.New("nil file")
}

func (fh *FileHandle) FormattedMTime() string {

	return fh.ModTime().Format(time.RFC822Z)
}

func (fh *FileHandle) GetContent() ([]byte, error) { //offset, limit int) (e error) {

	var buffer = make([]byte, fh.Size())
	fh.File.Seek(0, 0)
	len, err := fh.File.Read(buffer)
	if err != nil {
		return nil, err
	}

	if int64(len) != fh.Size() {
		return buffer, errors.New(fmt.Sprintf("read only %d of %d bytes", len, fh.Size()))
	}
	return buffer, nil
}

func (fh *FileHandle) SetContent(content []byte) error {

	return ioutil.WriteFile(fh.File.Name(), content, 0600)

	name := fh.File.Name()
	fd, err := os.OpenFile(name, os.O_RDWR, 0)
	if err != nil {
		return err
	}

	n_bytes, err := fd.WriteString(strings.Replace(string(content), "\r\n", "\n", -1))
	if err != nil {
		fmt.Println("error writing file:", err, n_bytes, fh.File)
	}
	return err
}

func (fh *FileHandle) Write(buffer []byte) (n int, err error) {
	n, err = fh.File.Write(buffer)
	return
}

func (fh *FileHandle) GetUTF8Content() (string, error) {

	content, err := fh.GetContent()
	if err != nil {
		return string(content), err
	}

	if utf8.Valid(content) {
		return string(content), nil
	} else {
		return string(content), errors.New("Invalid UTF-8")
	}
}

// ReadDir reads the directory named by dirname and returns
// a list of directory entries sorted by filename.
// from os.File.Readdir
func (i *FileHandle) ReadDir() ([]os.FileInfo, error) {

	list, err := i.File.Readdir(-1)
	if err != nil {
		return nil, err
	}
	//sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	return list, nil
}

func (fh *FileHandle) ReadDirHandle() ([]domain.Handle, error) {

	entries, err := fh.ReadDir()
	if err != nil {
		return nil, err
	}
	//sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })

	var handles = make([]domain.Handle, len(entries))
	for index, entry := range entries {
		if entry.Name()[0] == '.' {
			// ignore all files starting with '.'
			continue
		}
		handles[index] = &FileHandle{FileInfo: entry}
		//stat := entry.Sys()
		//.(*syscall.Stat_t) // _ ist ok und kein error
	}
	return handles, nil

	// file := FileFromInfo(info)

	// r, w, _ := file.GetPermissions(usr.Uid, usr.Gid)
	//file.Permissions = struct{ Read, Write bool }{Read: r, Write: w}
}

//PERMISSIONS

func (f *FileHandle) GetPermissions(uid, gid uint32) (r, w, x bool) { // => handle

	//fmt.Println(uid, gid)
	//Stat := f.Sys().(*syscall.Stat_t)

	//if Stat == nil {
	//	return false, false, false
	//}

	var fileUid uint32 //= -1 //Stat.Uid
	var fileGid uint32 //= -1 //Stat.Uid

	rwx := f.Mode().String()

	switch {
	case fileUid == uid:
		r, w, x = rwx[1] == 'r', rwx[2] == 'w', rwx[3] == 'x'
		// ur, uw, ux := f.UserPermissions()
		//r = r || ur
		//w = w || uw
		//x = x || ux
		//fmt.Println(r, w, x, ur, uw, ux)
		fallthrough
	case fileGid == gid:
		//gr, gw, gx := f.GroupPermissions()

		r = r || rwx[4] == 'r'
		w = w || rwx[5] == 'w'
		x = x || rwx[6] == 'x'
		//fmt.Println(r, w, x, gr, gw, gx)
		fallthrough
	default:
		//or, ow, ox := f.OthersPermissions()
		r = r || rwx[7] == 'r'
		w = w || rwx[8] == 'w'
		x = x || rwx[9] == 'x'
		//fmt.Println(r, w, x, or, ow, ox)
	}
	return r, w, x
}

func (f *FileHandle) UserPermissions() (r, w, x bool) {
	rwx := f.Mode().String()
	return rwx[1] == 'r', rwx[2] == 'w', rwx[3] == 'x'
}
func (f *FileHandle) GroupPermissions() (r, w, x bool) {
	rwx := f.Mode().String()
	return rwx[4] == 'r', rwx[5] == 'w', rwx[6] == 'x'
}
func (f *FileHandle) OthersPermissions() (r, w, x bool) {
	rwx := f.Mode().String()
	return rwx[7] == 'r', rwx[8] == 'w', rwx[9] == 'x'
}
