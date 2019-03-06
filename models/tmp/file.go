package models

import (
	"drive/auth"
	"drive/file"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/h2non/filetype/types"
)

type File struct {
	*storage.FileHandle
	Path        string `json:"path"`
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	Mode        os.FileMode
	ModTime     time.Time `json:"mtime"`
	AccessTime  time.Time
	ChangeTime  time.Time
	MIME        types.MIME
	Type        string `json:"type"`
	Owner       *user.User
	Group       *user.Group
	Permissions struct{ Read, Write bool }
}

func FileFromInfo(handle *storage.FileHandle) *File {

	file := &File{
		FileHandle: handle,
		//Path:       path,
		ModTime: handle.ModTime(),
		//AccessTime: statAtime(info.Stat),
		//ChangeTime: statCtime(info.Stat),
		Mode:  handle.Mode(),
		Name:  handle.Name(),
		Size:  handle.Size(),
		Owner: auth.GetUserByID(fmt.Sprintf("%d", handle.Stat.Uid)),
		Group: auth.GetGroupByID(fmt.Sprintf("%d", handle.Stat.Gid)),
	}
	//file.GuessMIME()
	return file
}

func (f *File) GetContent() ([]byte, error) { //offset, limit int) (e error) {

	var buffer = make([]byte, f.Size)
	f.Descriptor.Seek(0, 0)
	len, err := f.Descriptor.Read(buffer)
	if err != nil {
		return nil, err
	}

	if int64(len) != f.Size {
		return buffer, errors.New(fmt.Sprintf("read only %d of %d bytes", len, f.Size))
	}
	return buffer, nil
}

func (f *File) GetUTF8Content() (string, error) { //offset, limit int) (e error) {

	body, err := f.GetContent()
	if err != nil {
		return string(body), err
	}

	if utf8.Valid(body) {
		return string(body), nil
	} else {
		return string(body), errors.New("Invalid UTF-8")
	}

}

func (f *File) Save() error {
	//return ioutil.WriteFile(f.location, f.Content, 0600)
	return nil
}

func (f *File) SetContent(content []byte) error { //offset, limit int) (e error) {

	return ioutil.WriteFile(f.Descriptor.Name(), content, 0600)

	name := f.Descriptor.Name()
	fd, err := os.OpenFile(name, os.O_RDWR, 0)
	if err != nil {
		return err
	}

	n_bytes, err := fd.WriteString(strings.Replace(string(content), "\r\n", "\n", -1))
	if err != nil {
		fmt.Println("error writing file:", err, n_bytes, f.Descriptor)
	}
	return err
}

func (f *File) Breadcrumbs() []map[string]string {

	////////////////////

	var _ = func(path string) (dir, file string) {
		i := strings.LastIndex(path, "/")
		return path[:i+1], path[i+1:]
	}

	elements := strings.Split(strings.Trim(f.Path[1:], "/"), "/")
	breadcrumbs, currentPath := make([]map[string]string, len(elements)), ""
	for index, element := range elements {
		currentPath = currentPath + "/" + element
		breadcrumbs[index] = map[string]string{"name": element, "path": currentPath} // "/" + strings.Join(elements[:index+1], "/")}
	}
	breadcrumbs[len(elements)-1]["path"] = ""
	return breadcrumbs
}

func (f *File) BreadcrumbsAlt() []map[string]string {

	elements := strings.Split(strings.Trim(f.Path[1:], "/"), "/")
	breadcrumbs, currentPath := make([]map[string]string, len(elements)), ""
	for index, element := range elements {
		currentPath = currentPath + "/" + element
		breadcrumbs[index] = map[string]string{"name": element, "path": currentPath} // "/" + strings.Join(elements[:index+1], "/")}
	}
	breadcrumbs[len(elements)-1]["path"] = ""
	return breadcrumbs
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
func (f *File) ParentsAlt() []File {
	fmt.Println("path", f.Path)
	var path string
	elements := strings.Split(f.Path[1:], "/")
	fmt.Println("elements", elements)
	list := make([]File, len(elements))
	//fmt.Println("list", list)
	for index, element := range elements {
		path = fmt.Sprintf("%s/%s", path, element)
		list[index] = File{Name: element, Path: path}
		//fmt.Println(" - ", index, element)
	}
	//fmt.Println("list", list)
	return list
}

func (f *File) FormattedMTime() string {

	return f.ModTime.Format(time.RFC822Z)
}
func (f *File) String() string {

	return fmt.Sprintf("%s: %s", f.Type, f.Path)
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
	infos, err := ReadDir(parentPath)
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
