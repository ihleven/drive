package file

import (
	"drive/auth"
	"fmt"
	"os"
	"path"
	"syscall"
)

type Folder struct {
	File
	//Parent    string
	Folders   []*File
	Files     []*File
	IndexFile string
}

func NewDirectory(file *File, usr *auth.User) (*Folder, error) {
	file.Type = "D"
	dir := &Folder{File: *file}
	dir.List(usr)
	return dir, nil
}

func (d *Folder) List(usr *auth.User) error {

	entries, err := d.ReadDir()
	if err != nil {
		return err
	}
	for _, info := range entries {
		if info.Name()[0] == '.' {
			continue
		}
		file := d.NewChildFromFileInfo(info)
		r, w, _ := file.GetPermissions(usr.Uid, usr.Gid)
		file.Permissions = struct{ Read, Write bool }{Read: r, Write: w}

	}
	//d.Children = append(d.Folders, d.Files...)
	return nil
}
func (d *Folder) NewChildFromFileInfo(fileInfo os.FileInfo) *File {
	stat, _ := fileInfo.Sys().(*syscall.Stat_t) // _ ist ok und kein error

	info := &Info{FileInfo: fileInfo, Stat: stat}
	file, err := FileFromInfo(info)

	fmt.Println("NewChildFromFileInfo", err)
	file.Path = path.Join(d.Path, info.Name())
	if info.IsDir() {
		file.Type = "D"
		d.Folders = append(d.Folders, file)
	} else {
		//	file.GuessMIME()
		//child.ParseMIME()
		//child.MatchMIMEType()
		//child.DetectContentType()
		d.Files = append(d.Files, file)
	}
	return file
}
