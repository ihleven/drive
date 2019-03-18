package storage

import (
	"drive/domain"
	"errors"
	"os"
	"syscall"
	"time"
	"unicode/utf8"
)

type Meta struct {
	//Path string
	Stat *syscall.Stat_t

	Mode os.FileMode
	//Uid, Gid uint32
	//Size     int64
	//Created  time.Time
	//Modified time.Time
	//Accessed time.Time
	//MIME        types.MIME
	//Owner       uint32
	//Group       uint32
	Permissions struct{ Read, Write, Exec bool }
}

func (m *Meta) Name() string {
	return "this.is.the.name"
}
func (m *Meta) GetMode() os.FileMode {
	return m.Mode
}
func (m *Meta) Modified() time.Time {
	sec, nsec := m.Stat.Mtimespec.Unix()
	t := time.Unix(sec, nsec)
	return t
}

func (m *Meta) Size() int64 {
	return m.Stat.Size
}

func (m *Meta) Uid() uint32 {
	return m.Stat.Uid
}
func (m *Meta) Gid() uint32 {
	return m.Stat.Gid
}

func (m *Meta) HasReadPermission(uid, gid uint32) bool {

	if m.Mode&OS_OTH_R != 0 {
		return true
	}
	if m.Stat.Gid == gid && m.Mode&OS_GROUP_R != 0 {
		return true
	}
	if m.Stat.Uid == uid && m.Mode&OS_USER_R != 0 {
		return true
	}
	return false
}

func (m *Meta) ReadDirHandle() ([]domain.Handle, error) {
	return nil, nil
}
func (m *Meta) GetPermissions(uid, gid uint32) (r, w, x bool) {
	return false, false, false
}

func (m *Meta) GetContent() ([]byte, error) { //offset, limit int) (e error) {

	var buffer = make([]byte, m.Size())

	return buffer, nil
}

func (m *Meta) GetUTF8Content() (string, error) {

	content, err := m.GetContent()
	if err != nil {
		return string(content), err
	}

	if utf8.Valid(content) {
		return string(content), nil
	} else {
		return string(content), errors.New("Invalid UTF-8")
	}
}

func (m *Meta) SetContent(content []byte) error {

	return nil
}
