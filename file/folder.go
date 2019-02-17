package file

import (
	"os"
	"path"
)

type Folder struct {
	File
	//Parent    string
	Folders   []File
	Files     []File
	IndexFile string
}

func NewDirectory(file *File) (*Folder, error) {
	file.Type = "D"
	dir := &Folder{File: *file}
	dir.List()
	return dir, nil
}

func (d *Folder) List() error {

	entries, err := d.ReadDir()
	if err != nil {
		return err
	}
	for _, info := range entries {
		if info.Name()[0] == '.' {
			continue
		}
		d.NewChildFromFileInfo(info)

	}
	//d.Children = append(d.Folders, d.Files...)
	return nil
}
func (d *Folder) NewChildFromFileInfo(fileInfo os.FileInfo) *File {

	info := &Info{FileInfo: fileInfo}

	file := File{
		Info:     info,
		location: path.Join(d.location, info.Name()),
		Path:     path.Join(d.Path, info.Name()),
		Name:     info.Name(),
		Size:     info.Size(),
		Mode:     info.Mode(),
		ModTime:  info.ModTime(),
	}

	//file.GetPerm(&info)
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
	return &file
}
