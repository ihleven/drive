package domain

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/h2non/filetype/types"
)

type Account struct {
	Uid, Gid      uint32
	Username      string
	Name          string
	HomeDir       string
	Authenticated bool
}

type User struct {
	Uid      string
	Gid      string
	Username string
	Name     string
	HomeDir  string
}

type Group struct {
	Gid  string // group ID
	Name string // group name
}
type Mimetype struct {
	Type    string // group ID
	Subtype string // group name
	Charset string
}

type Storage interface {
	Open(string) (Handle, error)
	ReadDir(string) ([]Handle, error)
}

type Handle interface {
	os.FileInfo
	GuessMIME() types.MIME
	Close() error
	GetFile() *os.File
	ReadDir() ([]os.FileInfo, error)
	ReadDirHandle() ([]Handle, error)
	GetPermissions(uid, gid uint32) (r, w, x bool)
}

type File struct {
	Handle
	Path string `json:"path"`

	Name        string `json:"name"`
	Size        int64  `json:"size"`
	Mode        os.FileMode
	CTime       time.Time
	ATime       time.Time
	MTime       time.Time `json:"mtime"`
	Created     time.Time
	Modified    time.Time
	Accessed    time.Time
	MIME        types.MIME
	Owner       *User
	Group       *Group
	Permissions struct{ Read, Write, Exec bool }
}

func NewFile(handle Handle) *File {

	file := &File{
		Handle: handle,
		//Path:       path,
		MTime: handle.ModTime(),
		Owner: &User{},
		Group: &Group{},
		Mode:  handle.Mode(),
		Name:  handle.Name(),
		Size:  handle.Size(),
	}
	handle.GuessMIME()
	return file
}

func (f *File) Parents() []struct{ Name, Path string } {

	var path string
	elements := strings.Split(f.Path[1:], "/")
	list := make([]struct{ Name, Path string }, len(elements))
	for index, element := range elements {
		path = fmt.Sprintf("%s/%s", path, element)
		list[index] = struct{ Name, Path string }{Name: element, Path: path}
	}
	return list
}

type Folder struct {
	*File
	//Parent    string
	Entries   []*File
	IndexFile string
}

func NewDirectory(file *File, usr *Account) (*Folder, error) {

	//file.Type = "D"

	folder := &Folder{File: file}

	entries, err := file.Handle.ReadDirHandle()
	if err != nil {
		return nil, err
	}
	for _, info := range entries {
		if info.Name()[0] == '.' {
			continue
		}
		file := NewChildFromHandle(info, usr)
		folder.Entries = append(folder.Entries, file)
	}
	fmt.Println(folder.Entries)
	return folder, nil

}

func NewChildFromHandle(handle Handle, usr *Account) *File {

	file := NewFile(handle)

	//file.Path = path.Join(d.Path, info.Name())

	r, w, x := handle.GetPermissions(usr.Uid, usr.Gid)
	file.Permissions = struct{ Read, Write, Exec bool }{Read: r, Write: w, Exec: x}
	return file
}

type Album struct {
	// 	"github.com/eminetto/clean-architecture-go/pkg/entity"
	//ID   entity.ID `json:"id" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name,omitempty"`
	Path string `json:"path" bson:"path"`
}

type Repository interface {
	FindAll() ([]*Album, error)
	Get() (*Album, error)
}
