package drive

import (
	"drive/domain"
	"drive/errors"
	"encoding/json"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"path/filepath"
	"time"
)

// Album
//
// Album bettet *Directory und *File ein und wird über URLs wie "/alben/pfth/to/dir" erreicht.
//
// Metadaten in Datei meta.json => wird von struct AlbumMeta geparst.
type Album struct {
	*File `json:"file"`
	// Datei in der die folgenden Felder als JSON gespeichert werden.
	//AlbumFile   string `json:"metafile,omitempty"`
	Title       string   `json:"title"`
	Subtitle    string   `json:"subtitle"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
	// Die Bilder des Albums. Entweder alle Bilder des Verzeichnisses oder eine Liste von Bildern aus AlbumFile
	Images []Image2 `json:"images"`

	//Entries   []File `json:"entries"`
	Image string `json:"image"`
}

func GetAlbum(storage Storage, path string, usr *domain.Account) (*Album, error) {

	file, err := GetFile(storage, path, usr)
	if err != nil {
		return nil, errors.Wrap(err, "Could not get file for path %s", path)
	}
	if !file.HasReadPermission(usr.Uid, usr.Gid) {
		return nil, errors.New(403, "user: %v has no read permission for %v", usr, path)
	}

	album := Album{File: file}

	handles, err := storage.ReadDir(path)
	if err != nil {
		return nil, errors.Wrap(err, "Could not list album folder")
	}

	for _, handle := range handles {
		mime := handle.GuessMIME()
		switch {
		case mime.Type == "image":
			image, _ := NewImageFromHandle(handle)
			album.Images = append(album.Images, *image)
		case handle.Name() == "album.json":
			err := json.NewDecoder(handle.Descriptor(0)).Decode(&album)
			if err != nil {
				return nil, errors.Wrap(err, "Could not parse album file")
			}
		case file.MIME.Value == "text/diary; charset=utf-8":
			album.parseDiary(file.Name)
			//case strings.HasSuffix(file.Name, ".dia"):
			//	album.parseDiary(file.AsTextfile())
			//		album.parseMeta()
			//		album.Title = fmt.Sprintf("%s | alben", dir.Name)

		}

		path = filepath.Join(file.Path, handle.Name())

		//folder.Entries = append(folder.Entries, entry)

	}
	return &album, nil
}

func (a *Album) parseMeta() error {

	//content, err := ioutil.ReadFile(filepath.Join(a.location, "output.json"))

	//json.Unmarshal(content, a)

	return nil
}

func (a *Album) parseDiary(name string) (*Diary, error) {
	fmt.Println("parseDiary", name)
	return nil, nil
}

func (a *Album) Dump() error {
	_, err := json.MarshalIndent(a, "", "    ")
	if err != nil {
		return err
	}
	//err = ioutil.WriteFile(filepath.Join(a.location, "output.json"), data, 0644)
	//if err != nil {
	//	fmt.Println("ERROR dumping:", err)
	//	return err
	//}
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

/*
func NewDiary(fd *File) (*Diary, error) {

	input := bufio.NewScanner(fd.Descriptor)
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
			dir := filepath.Dir(fd.Path)
			img := filepath.Join(dir, strings.TrimSpace(line[2:]))
			file, err := file.Open(img)
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
*/
