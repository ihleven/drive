package drive

import (
	"drive/domain"
	"io"
	"os"
	"time"

	"github.com/h2non/filetype/types"
)

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
	Location(string) string
	GetHandle(string) (Handle, error)
	Open(string) (*os.File, error)
	Create(string) error
	Delete(string) error
	ReadDir(string) ([]Handle, error)
	Save(string, io.Reader) error
	//PermOpen(string, uint32, uint32) (*os.File, *time.Time, error)
	URL(string) string
	GetServeURL(string) string
	CleanPath(string) string
	CleanServePath(string) string
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
	Storage() Storage // liefert Storage, fuer die der Pfad gilt
	Location() string // Absoluter Pfad im Dateisystem
	Descriptor(int) *os.File

	StoragePath() string // Pfad ab storage root

	URL() string
	ServeURL() string
}

type Handle interface {
	os.FileInfo
	Locator

	ToFile(*domain.Account) (*File, error)
	GuessMIME() types.MIME

	//Read(b []byte) (n int, err error)
	//Write(p []byte) (n int, err error)

	GetContent() ([]byte, error)
	SetContent([]byte) error
	GetUTF8Content() (string, error)
	SetUTF8Content([]byte) error

	//Uid() uint32
	//Gid() uint32
	GetPermissions(owner uint32, group uint32, account *domain.Account) *Permissions
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
	Account *domain.Account `json:"-"`
	//Type              Permtype
	//Mode              os.FileMode
	IsOwner           bool
	InGroup           bool
	Read, Write, Exec bool
}

// File bundles all publically available information about Files (and Folders).
//
type File struct {
	Handle      `json:"-"`
	Name        string       `json:"name"`
	Path        string       `json:"path"`
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

type Folder struct {
	*File
	//Parent    string
	Entries   []*File `json:"entries"`
	IndexFile *File   `json:"indexFile"`
}
