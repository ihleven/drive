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

type File struct {
	Path string `json:"path"`

	Name        string `json:"name"`
	Size        int64  `json:"size"`
	Mode        os.FileMode
	MTime       time.Time `json:"mtime"`
	ATime       time.Time
	CTime       time.Time
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
