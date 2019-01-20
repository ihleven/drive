package fs

import (
	"bufio"
	"drive/gosrc/views"
	"encoding/json"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Album
//
// Album bettet *Directory und *File ein und wird über URLs wie "/alben/pfth/to/dir" erreicht.
//
// Metadaten in Datei meta.json => wird von struct AlbumMeta geparst.
type Album struct {
	*Directory `json:"-"`
	// Datei in der die folgenden Felder als JSON gespeichert werden.
	AlbumFile   string `json:"metafile,omitempty"`
	Title       string
	Description string
	Keywords    []string
	// Die Bilder des Albums. Entweder alle Bilder des Verzeichnisses oder eine Liste von Bildern aus AlbumFile
	Images []Image
}

func NewAlbum(dir *Directory) (*Album, error) {

	album := Album{Directory: dir}

	for i := 0; i < len(dir.Files); i++ {
		file := dir.Files[i]
		switch {
		case file.Name == "album.html":
			album.AlbumFile = file.Name

		//case strings.HasSuffix(file.Name, ".dia"):
		//	album.parseDiary(file.AsTextfile())

		case file.MIME.Value == "text/diary; charset=utf-8":
			album.parseDiary(file.Name)
		}

		if file.MIME.Type == "image" {
			img, err := file.AsImage()
			if err != nil {
				fmt.Println(" - ERROR '%s'\n\n\n", err)
			} else {
				album.Images = append(album.Images, *img)
			}
		}

	}
	album.parseMeta()
	album.Title = fmt.Sprintf("%s | alben", dir.Name)
	return &album, nil
}

func (a *Album) parseMeta() error {

	content, err := ioutil.ReadFile(filepath.Join(a.location, "output.json"))

	json.Unmarshal(content, a)

	return err
}

func (a *Album) parseDiary(name string) (*Diary, error) {
	fmt.Println("parseDiary", name)
	return nil, nil
}

func (a *Album) Dump() error {
	data, err := json.MarshalIndent(a, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(a.location, "output.json"), data, 0644)
	if err != nil {
		fmt.Println("ERROR dumping:", err)
		return err
	}
	return nil
}

func (a *Album) Render(w http.ResponseWriter, req *http.Request) error {

	contentType := req.Header.Get("Accept")

	if contentType == "application/json" {

		enc := json.NewEncoder(w)
		if err := enc.Encode(a); err != nil {
			log.Println(err)
		}
		// json, _ := json.Marshal(a)
		// w.Write(json)
	} else {
		fmt.Println("RENDER")
		err := views.Render("album", w, a)
		if err != nil {
			fmt.Println("ERROR:", err)
			//			panic(err)
		}
	}
	return nil
}

// DIARY
type Paragraph struct {
	Content string
}
type Diary struct {
	From       time.Time `json:"mtime"`
	IndexImage string
	Title      string
	Content    []Paragraph
	Images     []Image
}

func NewDiary(file *File, storage *FileSystemStorage) (*Diary, error) {
	fd, err := os.Open(file.location)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	input := bufio.NewScanner(fd)
	diary := &Diary{}
	currentParagraph := Paragraph{}
	for no := 1; input.Scan(); no++ {

		line := strings.TrimSpace(input.Text())
		switch {
		case strings.HasPrefix(line, "From:"):
			t, err := time.Parse("2006-01-02", line[5:])
			if err == nil {
				//diary.From = t
				fmt.Println(t)
			}
		case strings.HasPrefix(line, "I:"):
			dir := filepath.Dir(file.Path)
			img := filepath.Join(dir, strings.TrimSpace(line[2:]))
			file, err := storage.Open(img)
			//			root, e := filepath.Rel(file.location, file.Path)

			image, err := file.AsImage()
			if err != nil {
				return nil, err
			}
			diary.Images = append(diary.Images, *image)

		case line == "":
			diary.Content = append(diary.Content, currentParagraph)
			currentParagraph = Paragraph{}
		default:

			currentParagraph.Content += line
		}
	}
	return diary, nil
}

func (d *Diary) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := views.Render("diary", w, d)
	if err != nil {
		fmt.Println("ERROR:", err)
		//			panic(err)
	}
}
