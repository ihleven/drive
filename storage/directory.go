package storage

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
)

type Filer interface {
	RenderHTML(w http.ResponseWriter, req *http.Request)
	MarschalJSON(w http.ResponseWriter, req *http.Request)
}

type Directory struct {
	Dirs      []string
	Files     []string
	IndexFile string
}
type Dirlisting struct {
	Name           string
	Children_dir   []string
	Children_files []string
	ServerUA       string
}

func (d *Directory) RenderHTML(w http.ResponseWriter, req *http.Request) {

	tpl, err := template.ParseFiles("templates/directory.html")
	if err != nil {
		http.Error(w, "500 Internal Error : Error while generating directory listing.", 500)
		fmt.Println(err)
		log.Print(err)
		return
	}

	data := Dirlisting{Name: req.URL.Path, ServerUA: "ihle",
		Children_dir: d.Dirs, Children_files: d.Files}

	err = tpl.Execute(w, data)
	if err != nil {
		fmt.Println(err)
	}
}

func (d *Directory) MarschalJSON(w http.ResponseWriter, req *http.Request) {

}

func NewDirectory(fileInfo os.FileInfo, dirname string) (*Directory, error) {

	//fileInfos, err := ioutil.ReadDir(dirname)
	//if err != nil {
	//	log.Fatal(err)
	//}
	// oder ...
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	fileInfos, _ := f.Readdir(-1) // returns all the FileInfo from the directory in a single slice
	f.Close()
	if err != nil {
		return nil, err
	}
	// sorting, TODO: dirs first...
	sort.Slice(fileInfos, func(i, j int) bool { return fileInfos[i].Name() < fileInfos[j].Name() })
	//return list, nil
	// end oder

	d := &Directory{Dirs: make([]string, 0), Files: make([]string, 0)}

	for _, val := range fileInfos {
		if val.Name()[0] == '.' {
			continue
		} // Remove hidden files from listing
		if val.Name() == "index.html" {
			//handleFile(path.Join(f.Name(), "index.html"), w, req)
			d.IndexFile = val.Name()
		}
		mode := val.Mode()

		if mode.IsDir() {
			d.Dirs = append(d.Dirs, val.Name())

		} else if mode.IsRegular() {
			d.Files = append(d.Files, val.Name())
		}
	}
	return d, nil
}
