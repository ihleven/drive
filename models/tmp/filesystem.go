package models

import "drive/file"

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

type FS interface {
	GetFileLikeObject(filesystemName, username, path string) (*file.File, error)
}

type Filesystem struct {
	Name string // Name des Filesystems
	Root string // Pfad zum Root des Filesystems
	User string // Angemeldeter User, dem die die Dateien gehören
}

func (fs *Filesystem) Open(path string) (*file.File, error) {
	return nil, nil
}
