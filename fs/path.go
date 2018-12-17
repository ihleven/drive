package fs

import (
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
)

func NewPath(root, pathname string) (*Path, error) {

	fullpath := path.Join(root, pathname)

	fileInfo, error := os.Stat(fullpath)
	if error != nil {
		if os.IsNotExist(error) {
			return nil, nil
		}
		return nil, error
	}
	fp := Path{
		Abs:     fullpath,
		Root:    &root,
		Path:    pathname,
		Name:    fileInfo.Name(),
		Size:    fileInfo.Size(),
		Mode:    fileInfo.Mode(),
		ModTime: fileInfo.ModTime()}

	// , Abs: fullpath Dirname: dirname, Filename: filename}
	fp.DetectContentType()
	fp.CalcParents()
	return &fp, nil
}

type Path struct {
	Root *string // /home/ihle
	Path string  // /alben/urlaube/Wien2013/IMG_123.jpg
	Abs  string  // /home/ihle/alben/urlaube/Wien2013/IMG_123.jpg
	//Dirname  string // /alben/urlaube/Wien2013/
	//Filename string // IMG_123.jpg
	Name     string
	Size     int64
	Mode     os.FileMode
	ModTime  time.Time
	MIME     types.MIME
	Children []Path
	Parents  []Path
}

func (p *Path) CalcParents() {
	var pa string
	elements := strings.Split(p.Path[1:], "/")
	list := make([]Path, len(elements))
	for index, element := range elements {
		pa = fmt.Sprintf("%s/%s", pa, element)
		list[index] = Path{Name: element, Path: pa}
		fmt.Println(index, list[index].Name, list[index].Path)
	}
	p.Parents = list
}

func (p *Path) ParseMIME() {

	if ext := path.Ext(p.Name); ext != "" {
		mimetype := mime.TypeByExtension(ext)
		p.MIME = filetype.GetType(ext[1:]).MIME //   types.Get(ext).MIME
		fmt.Printf(" * MIME '%s' => %s, %s, %s\n", p.Name, mimetype, p.MIME.Value, mime.TypeByExtension(ext))
	}
}

func (f *Path) MatchMIMEType() error {

	file, err := os.Open(f.Abs)
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
func (f *Path) DetectContentType() error {

	fd, err := os.Open(f.Abs)
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
func (p *Path) Type() string {
	if p.IsDir() {
		return "DIR"
	} else if p.Mode.IsRegular() {
		return "FILE"
	} else {
		return ""
	}
}

func (p *Path) IsDir() bool     { return p.Mode.IsDir() }
func (p *Path) IsRegular() bool { return p.Mode.IsRegular() }

func (p *Path) String() string {
	return fmt.Sprintf("%s: %s", p.Type(), p.Path)
}

func (p *Path) NewChild(info os.FileInfo) Path {

	return Path{
		Abs:     path.Join(p.Abs, info.Name()),
		Root:    p.Root,
		Path:    path.Join(p.Path, info.Name()),
		Name:    info.Name(),
		Size:    info.Size(),
		Mode:    info.Mode(),
		ModTime: info.ModTime()}
}

func (p *Path) List() error {

	fileInfos, err := ioutil.ReadDir(p.Abs)
	if err != nil {
		log.Fatal(err)
		return err
	}
	for _, info := range fileInfos {
		if info.Name()[0] == '.' {
			continue
		}
		child := p.NewChild(info)
		child.ParseMIME()
		child.MatchMIMEType()
		child.DetectContentType()
		p.Children = append(p.Children, child)
	}
	return nil
}
