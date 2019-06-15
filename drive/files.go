package drive

import (
	"drive/domain"
	"drive/errors"
	"fmt"
	"path"
	"path/filepath"
	"strings"
)

// GetReadHandle bekommt einen bereinigten Pfad
func GetReadHandle(storage Storage, path string, uid, gid uint32) (Handle, error) {

	handle, err := storage.GetHandle(path)
	if err != nil {
		return nil, errors.Wrap(err, "Could not get file handle")
	}
	if !handle.HasReadPermission(uid, gid) {
		return nil, errors.New(403, "uid: %v, gid %v has not read permission for %v", uid, gid, path)
	}
	return handle, nil
}

func GetFile(storage Storage, path string, usr *domain.Account) (*File, error) {
	fmt.Println("GetFile", path)

	handle, err := storage.GetHandle(storage.CleanPath(path))
	if err != nil {
		return nil, errors.Wrap(err, "Could not get file handle for %s", path)
	}

	file, err := handle.ToFile(usr)
	if err != nil {
		return nil, errors.Wrap(err, "Could not transform handle %v to File", file)
	}
	fmt.Println("file:", file)
	return file, nil
}

func CreateFile(storage Storage, path string, usr *domain.Account) (*File, error) {

	err := storage.Create(path)
	if err != nil {
		return nil, errors.Wrap(err, "File '%s' could not be created!", path)
	}
	return GetFile(storage, path, usr)
}

func DeleteFile(file *File) error {

	if !file.Permissions.Write {
		return errors.New(errors.PermissionDenied, "Missing write permission for %s", file.Path)
	}
	err := file.Storage().Delete(file.Path)
	if err != nil {
		return errors.Wrap(err, "File '%s' could not be deleted!", file.Path)
	}
	return nil
}

func GetFolder(file *File, usr *domain.Account) (*Folder, error) {

	folder := &Folder{File: file}
	handles, err := file.Storage().ReadDir(file.StoragePath())
	if err != nil {
		return nil, err
	}

	for _, handle := range handles {

		entry, _ := handle.ToFile(usr) // filepath.Join(file.Path, handle.Name()),

		folder.Entries = append(folder.Entries, entry)
		if entry.Name == "index.html" {
			folder.IndexFile = entry
		}
	}
	return folder, nil
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
	parent := path.Dir(f.StoragePath())
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
		return nil, errors.Wrap(err, "Could not read dir %s", parentPath)
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
