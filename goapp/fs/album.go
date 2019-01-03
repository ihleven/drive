package fs

import (
	"drive/goapp/views"
	"encoding/json"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
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

type Diary struct {
	*Album
}

func NewAlbum(dir *Directory) (*Album, error) {

	album := Album{Directory: dir}

	//for _, file := range dir.Files {
	for i := 0; i < len(dir.Files); i++ {
		file := dir.Files[i]
		if file.Name == "album.html" {
			album.AlbumFile = file.Name
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

		err := views.Render("album", w, a)
		if err != nil {
			fmt.Println("ERROR:", err)
			//			panic(err)
		}
	}
	return nil
}
