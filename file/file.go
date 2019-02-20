package file

import (
	"drive/auth"
	"errors"
	"fmt"
	"os"
	"os/user"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/h2non/filetype/types"
)

type File struct {
	*Info
	location    string
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

func NewFileLöschen(path string, usr *auth.User) (*File, error) {
	fmt.Println("NewFile", path, usr)
	info, err := Open(path, 0, 0, 0)
	if err != nil {
		return nil, err
	}

	f, err := FileFromInfo(info)
	r, w, _ := f.GetPermissions(usr.Uid, usr.Gid)
	f.Permissions = struct{ Read, Write bool }{Read: r, Write: w}

	return f, nil
}

func FileFromInfo(info *Info) (*File, error) {
	fmt.Println("info.Stat.Uid", info.Stat.Uid)
	owner, err := user.LookupId(fmt.Sprintf("%d", info.Stat.Uid))
	if err != nil {
		fmt.Println("err", err)
		switch err.(type) {
		case user.UnknownUserIdError:
			owner = &user.User{Username: "unknown"}
		default:
			owner = &user.User{Username: "unknown"}
		}
	}
	group, err := user.LookupGroupId(fmt.Sprintf("%d", info.Stat.Gid))
	if err != nil {
		switch err.(type) {
		case user.UnknownGroupIdError:
			group = &user.Group{Name: "unknown"}
		default:
			group = &user.Group{Name: "unknown"}
		}
	}

	f := &File{
		Info: info,
		//Path:       path,
		ModTime:    info.ModTime(),
		AccessTime: statAtime(info.Stat),
		ChangeTime: statCtime(info.Stat),
		Mode:       info.Mode(),
		Name:       info.Name(),
		Size:       info.Size(),
		Owner:      owner,
		Group:      group,
	}
	fmt.Println("size", f.Size)
	return f, nil
}

func (f *File) GetContent() ([]byte, error) { //offset, limit int) (e error) {

	var buffer = make([]byte, f.Size)
	len, err := f.Descriptor.Read(buffer)
	if err != nil {
		return nil, err
	}

	if int64(len) != f.Size {
		return buffer, errors.New(fmt.Sprintf("read only %d of %d bytes", len, f.Size))
	}
	return buffer, nil
}

func (f *File) StringContent() string {
	c, _ := f.GetContent()
	return string(c)
}

func (f *File) GetTextContent() (string, error) { //offset, limit int) (e error) {
	var body = make([]byte, f.Info.Size())
	fmt.Println("GetTextContent", f.Info.Size())

	len, err := f.Descriptor.Read(body)
	if err != nil {
		return "", err
	}
	if int64(len) != f.Size {
		fmt.Println("GetTextContent", len, f.Info.Size())
		return "", errors.New(fmt.Sprintf("read only %d of %d bytes", len, f.Size))
	}
	if utf8.Valid(body) {
		return string(body), nil
	} else {
		return "", errors.New("Invalid UTF-8")
	}

}

func (f *File) Save() error {
	//return ioutil.WriteFile(f.location, f.Content, 0600)
	return nil
}

func (f *File) SetContent(content string) error { //offset, limit int) (e error) {

	//body, e := ioutil.ReadFile(f.location)
	//if utf8.Valid(body) {
	//	f.Content = body
	//} else {
	//	e = errors.New("Invalid UTF-8")
	//}
	return nil
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
