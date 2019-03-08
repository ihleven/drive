package file

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func (f *File) Save() error {
	//return ioutil.WriteFile(f.location, f.Content, 0600)
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
