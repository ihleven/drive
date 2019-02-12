package models

import (
	"os"
	"time"

	"github.com/h2non/filetype/types"
)

/*
Files:
======
https://ihle.fm/home/go/src/main.go							=>	homeHandler(user, "go/src/main.go")
https://ihle.fm/home/go/src/main.go
https://ihle.fm/home/go/src/main.go
https://ihle.fm/serve/~matt/go/src/main.go


https://ihle.fm/fotos/urlaube/src/main.go
https://ihle.fm/public/weihnachtsfeier2018/IMG_5673.jpg
https://ihle.fm/home/weihnachtsfeier2018/IMG_5673.jpg




*/

type File struct {
	location string
	Path     string `json:"path"`
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Mode     os.FileMode
	ModTime  time.Time `json:"mtime"`
	MIME     types.MIME
	Type     string `json:"type"`
}

type FS interface {
	GetFileLikeObject(filesystemName, username, path string) (*File, error)
}

type Filesystem struct {
	Name string // Name des Filesystems
	Root string // Pfad zum Root des Filesystems
	User string // Angemeldeter User, dem die die Dateien gehören
}

func (fs *Filesystem) Open(path string) (*File, error) {
	return fs.GetFileLikeObject()
}
