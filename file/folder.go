package file

import (
	"drive/auth"
	"os"
	"path"
	"syscall"
)

type Folder struct {
	File
	//Parent    string
	Folders   []*File
	Files     []*File
	Entries   []*File
	IndexFile string
}

func NewDirectory(file *File, usr *auth.Account) (*Folder, error) {

	file.Type = "D"

	dir := &Folder{File: *file}

	entries, err := dir.ReadDir()
	if err != nil {
		return nil, err
	}

	for _, info := range entries {
		if info.Name()[0] == '.' {
			continue
		}
		file := dir.NewChildFromFileInfo(info, usr)
		dir.Entries = append(dir.Entries, file)
	}
	return dir, nil

}
func (d *Folder) NewChildFromFileInfo(fileInfo os.FileInfo, usr *auth.Account) *File {
	stat, _ := fileInfo.Sys().(*syscall.Stat_t) // _ ist ok und kein error

	info := &Info{FileInfo: fileInfo, Stat: stat}
	file := FileFromInfo(info)

	file.Path = path.Join(d.Path, info.Name())

	r, w, _ := file.GetPermissions(usr.Uid, usr.Gid)
	file.Permissions = struct{ Read, Write bool }{Read: r, Write: w}
	return file
}
