package drive

import (
	"fmt"
	"io"
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
	Uid      string `json:"uid"`
	Gid      string `json:"-"`
	Username string `json:"name"`
	Name     string `json:"-"`
	HomeDir  string `json:"-"`
}

type Group struct {
	Gid  string `json:"gid"`  // group ID
	Name string `json:"name"` // group name
}
type Mimetype struct {
	Type    string `json:"type"`
	Subtype string `json:"subtype"`
	Charset string `json:"charset"`
}

type Storage interface {
	GetHandle(name string) (Handle, error)
	Open(string) (*os.File, error)
	Create(string) error
	Delete(string) error
	ReadDir(string) ([]Handle, error)
	Save(string, io.Reader) error

	//PermOpen(string, uint32, uint32) (*os.File, *time.Time, error)
}

type Mimer interface {
	Guess() (Mimetype, error)
	MIME() string
	Parts() (string, string, string)
	Type() string
	Subtype() string
	Charset() string
}
type Locator interface {
	Storage() Storage
	Location() string
	Descriptor(int) *os.File
}

type Handle interface {
	os.FileInfo
	Storage() Storage
	Descriptor(int) *os.File

	ToFile(string, *Account) (*File, error)
	GuessMIME() types.MIME

	ListDirHandles(bool) ([]Handle, error)

	GetContent() ([]byte, error)
	SetContent([]byte) error
	GetUTF8Content() (string, error)
	SetUTF8Content([]byte) error

	//Uid() uint32
	//Gid() uint32
	GetPermissions(owner uint32, group uint32, account *Account) *Permissions
	HasReadPermission(uid, gid uint32) bool
}

type UTF8er interface {
	GetUTF8Content() (string, error)
	SetUTF8Content([]byte) error
}
type Loader interface {
	Load() ([]byte, error)
}
type Saver interface {
	Save([]byte) error
}
type Loaver interface {
	Load() ([]byte, error)
	Save([]byte) error
	Write(p []byte) (n int, err error)
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
	Handle      `json:"-"`
	Path        string       `json:"path"`
	Name        string       `json:"name"`
	Size        int64        `json:"size"`
	Mode        os.FileMode  `json:"mode"`
	Owner       *User        `json:"owner"`
	Group       *Group       `json:"group"`
	Permissions *Permissions `json:"permissions"`
	Created     *time.Time   `json:"created"`
	Modified    time.Time    `json:"modified"`
	Accessed    *time.Time   `json:"accessed"`
	MIME        types.MIME   `json:"mime"`
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

func (f *File) Siblings() (*Siblings, error) {
	var currentIndex int
	siblings := &Siblings{}
	parentPath := f.ParentPath()
	infos, err := f.Storage().ReadDir(parentPath)
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
	Entries   []*File `json:"entries"`
	IndexFile *File   `json:"indexFile"`
}
