package domain

import (
	"os"
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
}

type Handle interface {
	os.FileInfo
	GuessMIME() types.MIME
	Close() error
	GetFile() *os.File
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

type Folder struct {
	*File
	//Parent    string
	Entries   []*File
	IndexFile string
}
