package fs

import (
	"drive/views"
	"fmt"
	_ "image/jpeg"
	"mime"
	"net/http"
	"os"
	"os/user"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
)

type FileHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}
type File struct {
	location string
	Path     string `json:"path"`
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Mode     os.FileMode
	ModTime  time.Time `json:"mtime"`
	MIME     types.MIME
	Type     string `json:"type"`
	User     *user.User
	Group    *user.Group
}

func NewFile(name string, fullpath string, fi *os.FileInfo) *File {

	//fmt.Printf("%s %s %s %s", uid, gid, Group.Name, User.Username)

	f := &File{
		location: fullpath,
		Path:     name,
		Name:     (*fi).Name(),
		Size:     (*fi).Size(),
		Mode:     (*fi).Mode(),
		ModTime:  (*fi).ModTime(),
	}
	f.GetPerm(fi)
	return f
}

func (f *File) GetPerm(fi *os.FileInfo) error {

	stat := (*fi).Sys().(*syscall.Stat_t)
	uid := fmt.Sprintf("%d", stat.Uid)
	gid := fmt.Sprintf("%d", stat.Gid)
	f.Group, _ = user.LookupGroupId(gid)
	f.User, _ = user.LookupId(uid)
	return nil
}

func (f *File) Breadcrumbs() []map[string]string {

	elements := strings.Split(strings.Trim(f.Path[1:], "/"), "/")
	breadcrumbs, currentPath := make([]map[string]string, len(elements)), ""
	for index, element := range elements {
		currentPath = currentPath + "/" + element
		breadcrumbs[index] = map[string]string{"name": element, "path": currentPath} // "/" + strings.Join(elements[:index+1], "/")}
	}
	breadcrumbs[len(elements)-1]["path"] = ""
	return breadcrumbs
}

func (f *File) Parents() []File {
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

func (f *File) String() string {

	return fmt.Sprintf("%s: %s", f.Type, f.Path)
}

func (f *File) FormattedMTime() string {

	return f.ModTime.Format(time.RFC822Z)
}

// func (f *File) GetTypeUnused() string {
// 	if f.IsDir() {
// 		return "DIR"
// 	} else if f.Mode.IsRegular() {
// 		return "FILE"
// 	} else {
// 		return ""
// 	}
// }

func (f *File) IsDir() bool {

	return (*f).Mode.IsDir()
}
func (f *File) IsRegular() bool { return f.Mode.IsRegular() }

func (f *File) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	views.SerializeJSON(w, f)
}

func (f *File) Specific() (fh FileHandler, err error) {

	if f.IsDir() {
		fh, err = NewDirectory(f)
		return
	}

	if f.IsRegular() {
		f.Type = "F"
		f.GuessMIME()
		switch f.MIME.Type {
		case "image":
			fh, err = f.AsImage()
			//image.ServeHTTP(w, r)

		case "text":
			switch f.MIME.Subtype {
			case "diary; charset=utf-8":
				//fh, err = NewDiary(f)
				fh, err = f.AsTextfile()
				fmt.Println("DIARY", fh)
			default:
				fh, err = f.AsTextfile()
				fmt.Println("TESTFILE", fh, f.MIME.Subtype)
			}

		default:
			fmt.Println("No specific file type found")
			fh, err = f, nil
		}

	}
	return
}

func (f *File) GuessMIME() {

	if ext := path.Ext(f.Name); ext != "" {
		mime := mime.TypeByExtension(ext) // mime package
		if mime != "" {
			f.MIME = types.NewMIME(mime)

		} else {
			f.MIME = filetype.GetType(ext[1:]).MIME // h2non/filetype => types.Get(ext).MIME
			//if f.MIME.Value != mimetype {
			//	fmt.Printf("MIME (%s): %s != %s \n", f.Name, mimetype, f.MIME)
			//}

		}
	}

	if f.MIME.Value == "" {
		f.h2nonMatchMIME261()
	}

	switch f.MIME.Type {
	case "image":
		f.Type = "FI"
	case "text":
		f.Type = "FT"
	default:
		f.Type = "F"
	}

}
func (f *File) h2nonMatchMIME261() error {

	file, err := os.Open(f.location)
	if err != nil {
		return err
	}
	defer file.Close()

	// https://github.com/h2non/filetype
	// We only have to pass the file header = first 261 bytes
	head := make([]byte, 261)
	file.Read(head)
	if t, e := filetype.Match(head); e == nil {
		f.MIME = t.MIME
	}
	fmt.Printf(" * MIME2 '%s' => %s, %s, %s\n", f.Name, f.MIME.Value, f.MIME.Type, f.MIME.Subtype)

	return nil
}

func (f *File) DetectContentType() error {

	fd, err := os.Open(f.location)
	if err != nil {
		return err
	}
	defer fd.Close()

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err = fd.Read(buffer)
	if err != nil {
		return err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)
	f.MIME = types.NewMIME(contentType)
	fmt.Printf(" * http '%s' => %s\n", f.Name, f.MIME.Type)

	return nil
}
