package domain

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
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
	GetHandle(name string) (Handle, error)
	//Open(string) (Handle, error)
	ReadDir(string) ([]Handle, error)
	//PermOpen(string, uint32, uint32) (*os.File, *time.Time, error)
	OpenFD(name string) (*os.File, error)
}

type Handle interface {
	os.FileInfo
	Descriptor() *os.File
	ToFile(string, *Account) (*File, error)
	GuessMIME() types.MIME
	//Close() error
	//GetFile() *os.File
	//ReadDir() ([]os.FileInfo, error)
	ReadDirHandle() ([]Handle, error)
	GetPermissions(owner uint32, group uint32, account *Account) (*Permissions, error)
	GetContent() ([]byte, error)
	SetContent([]byte) error
	GetUTF8Content() (string, error)
	//Uid() uint32
	//Gid() uint32
	HasReadPermission(uid, gid uint32) bool
}
type Permtype int

const (
	UNKNOWN Permtype = iota
	PERM_USR
	PERM_GRP
	PERM_OTH
)

type Permissions struct {
	Account *Account
	//Type              Permtype
	//Mode              os.FileMode
	IsOwner           bool
	InGroup           bool
	Read, Write, Exec bool
}

type File struct {
	Handle
	Path        string `json:"path"`
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	Mode        os.FileMode
	Owner       *User
	Group       *Group
	Permissions *Permissions
	Created     *time.Time
	Modified    time.Time
	Accessed    *time.Time
	MIME        types.MIME
}

func NewFile(handle Handle, path string) *File {

	file := &File{
		Handle: handle,
		Path:   path,
		Name:   handle.Name(),
		//Size:   handle.Size(),
		//Mode:   handle.Mode(),
		//MTime:  handle.ModTime(),
		//MIME:   handle.GuessMIME(),
		Owner: &User{},
		Group: &Group{},
	}
	return file
}

func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
func (f *File) Breadcrumbs() (breadcrumbs []struct{ Name, Path string }) {
	var url = "/"
	for elem, remainder := ShiftPath(f.Path); elem != ""; elem, remainder = ShiftPath(remainder) {
		url = path.Join(url, elem)
		breadcrumbs = append(breadcrumbs, struct{ Name, Path string }{Name: elem, Path: url})
	}
	return
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

func (f *File) ParentsWithFiles() []File {
	var path string
	elements := strings.Split(f.Path[1:], "/")
	list := make([]File, len(elements))
	for index, element := range elements {
		path = fmt.Sprintf("%s/%s", path, element)
		list[index] = File{Name: element, Path: path}
	}
	return list
}

func (f *File) ParentPath() string {
	parent := path.Dir(f.Path)
	if parent == "." {
		return ""
	}
	return parent
}

type Siblings struct {
	Count, CurrentIndex, PrevIndex, NextIndex int
	First,
	Last,
	Prev,
	Current,
	Next string
	All []string
}

func (f *File) Siblings(storage Storage) (*Siblings, error) {
	var currentIndex int
	siblings := &Siblings{}
	parentPath := f.ParentPath()
	infos, err := storage.ReadDir(parentPath)
	if err != nil {
		return nil, err
	}
	for _, info := range infos {

		if info.Name()[0] == '.' || info.IsDir() || filepath.Ext(info.Name()) != ".jpg" {
			continue
		}
		currentIndex++
		if info.Name() == f.Name {
			siblings.Current = path.Join(parentPath, info.Name())
			siblings.CurrentIndex = currentIndex
		}
		siblings.All = append(siblings.All, path.Join(parentPath, info.Name()))
	}

	siblings.Count = len(siblings.All)
	if siblings.Count > 0 {
		siblings.First = siblings.All[0]
		siblings.Last = siblings.All[siblings.Count-1]
		if siblings.CurrentIndex > 1 {
			siblings.Prev = siblings.All[siblings.CurrentIndex-2]
			siblings.PrevIndex = siblings.CurrentIndex - 1
		}
		if siblings.CurrentIndex < siblings.Count {
			siblings.Next = siblings.All[siblings.CurrentIndex]
			siblings.NextIndex = siblings.CurrentIndex + 1
		}

	}
	return siblings, nil
}

type Folder struct {
	*File
	//Parent    string
	Entries   []*File
	IndexFile *File
}

func NewDirectory(file *File, usr *Account) (*Folder, error) {

	//file.Type = "D"
	fmt.Println("NEwDirectory:", file.Path)
	folder := &Folder{File: file}

	entries, err := file.Handle.ReadDirHandle()
	if err != nil {
		return nil, err
	}
	for _, info := range entries {
		if info.Name()[0] == '.' {
			continue
		}
		f := NewFile(info, filepath.Join(file.Path, info.Name()))

		//f := NewChildFromHandle(info, file.Path, usr)
		folder.Entries = append(folder.Entries, f)
	}
	return folder, nil

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
