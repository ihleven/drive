package storage

import (
	"drive/domain"
	"drive/domain/usecase"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"syscall"
	"unicode/utf8"
)

type FileHandle struct {
	os.FileInfo
	storage  domain.Storage
	location string
	mode     os.FileMode
}

func NewFileHandle(info os.FileInfo, st *FileSystemStorage, location string) *FileHandle {

	handle := &FileHandle{
		FileInfo: info,
		storage:  st,
		mode:     info.Mode(),
		location: location,
	}
	if st.PermissionMode != 0 {
		handle.mode = (handle.mode & 0xfffffe00) | (st.PermissionMode & 0x1ff)
	}
	return handle
}

func (fh *FileHandle) ToFile(path string, account *domain.Account) (*domain.File, error) {

	stat, err := fh.Stat()
	if err != nil {
		return nil, err
	}
	permissions, err := fh.GetPermissions(stat.Uid, stat.Gid, account)
	if err != nil {
		log.Fatal("\nToFile:", err.Error())
		return nil, nil
	}
	file := &domain.File{
		Handle:      fh,
		Path:        path,
		Name:        fh.Name(),
		Size:        fh.Size(),
		Mode:        fh.mode,
		Owner:       usecase.GetUserByID(stat.Uid),
		Group:       usecase.GetGroupByID(stat.Gid),
		Permissions: permissions,
		Created:     statCtime(stat),
		Modified:    fh.ModTime(),
		Accessed:    statAtime(stat),
		MIME:        fh.GuessMIME(),
	}

	return file, nil
}

func (fh *FileHandle) Descriptor(flag int) *os.File { // , perm os.FileMode

	fd, err := os.OpenFile(fh.location, flag, 0755)
	if err != nil {
		log.Fatal("error gettting descriptor", err.Error(), fh.location)
		return nil
	}
	return fd
}

func (fh *FileHandle) Stat() (*syscall.Stat_t, error) {
	stat, ok := fh.Sys().(*syscall.Stat_t)
	if !ok {
		return nil, errors.New("Sys().(*syscall.Stat_t) => underlying data source is nil")
	}
	return stat, nil
}

func (fh *FileHandle) HasReadPermission(uid, gid uint32) bool {

	if fh.mode&OS_OTH_R != 0 {
		return true
	}
	//owner, group, err := fh.GetOwnerAndGroupID()
	stat, ok := fh.Sys().(*syscall.Stat_t)
	if !ok {
		log.Fatal(errors.New("Sys().(*syscall.Stat_t) => underlying data source is nil"))
		return false
	}
	if stat.Gid == gid && fh.mode&OS_GROUP_R != 0 {
		return true
	}
	if stat.Uid == uid && fh.mode&OS_USER_R != 0 {
		return true
	}
	return false
}

//PERMISSIONS

func (fh *FileHandle) GetPermissions(owner uint32, group uint32, account *domain.Account) (*domain.Permissions, error) { // => handle

	perm := &domain.Permissions{IsOwner: account.Uid == owner, InGroup: account.Gid == group}

	rr, wr, xr := OS_OTH_R, OS_OTH_W, OS_OTH_X
	if perm.InGroup {
		rr, wr, xr = rr|OS_GROUP_R, wr|OS_GROUP_W, xr|OS_GROUP_X
	}
	if perm.IsOwner {
		rr, wr, xr = rr|OS_USER_R, wr|OS_USER_W, xr|OS_USER_X
	}

	perm.Read = int(fh.mode)&rr != 0
	perm.Write = int(fh.mode)&wr != 0
	perm.Exec = int(fh.mode)&xr != 0
	return perm, nil
}

func (fh *FileHandle) ListDirHandles(hideDotFiles bool) ([]domain.Handle, error) {

	fd, err := os.Open(fh.location)
	if err != nil {
		return nil, err
	}
	entries, err := fd.Readdir(-1)
	fd.Close()
	if err != nil {
		return nil, err
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })

	var handles = make([]domain.Handle, 0)
	for _, entry := range entries {
		if hideDotFiles && entry.Name()[0] == '.' {
			// ignore all files starting with '.'
			continue
		}

		handles = append(handles, NewFileHandle(
			entry,
			fh.storage.(*FileSystemStorage),
			filepath.Join(fh.location, entry.Name()),
		))
	}
	return handles, nil
}

func (fh *FileHandle) Storage() domain.Storage {
	return fh.storage
}

/////////////

func (fh *FileHandle) GetContent() ([]byte, error) { //offset, limit int) (e error) {

	var content = make([]byte, fh.Size())

	fd := fh.Descriptor(0)
	defer fd.Close()
	fd.Seek(0, 0)

	bytes, err := fd.Read(content)
	if err != nil {
		return nil, err
	}

	if int64(bytes) != fh.Size() {
		return content, errors.New(fmt.Sprintf("read only %d of %d bytes", bytes, fh.Size()))
	}
	return content, nil
}

func (fh *FileHandle) SetContent(content []byte) error {

	fd := fh.Descriptor(os.O_RDWR | os.O_CREATE | os.O_TRUNC)

	n_bytes, err := fd.WriteAt(content, 0)
	if err != nil {
		fmt.Println("error writing file:", err, n_bytes)
		return err
	}
	fd.Close()
	return nil
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

func (fh *FileHandle) SetUTF8Content(content []byte) error {

	if !utf8.Valid(content) {
		return errors.New("Invalid UTF-8")
	}
	return fh.SetContent(content)
}

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
