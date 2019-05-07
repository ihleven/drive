package drive

import (
	"bufio"
	"drive/domain"
	"drive/errors"
	"encoding/json"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
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
	Image   string   `json:"image"`
	From    string   `json:"from"`
	Until   string   `json:"until"`
	Diaries []Diary  `json:"diaries"`
	Moments []Album  `json:"moments"`
	Sources []Source `json:"sources"`
	Pages   []string `json:"pages"`
}
type Source struct {
	Name         string   `json:"name"`
	Camera       string   `json:"camera"`
	Photographer string   `json:"photographer"`
	Images       []Image2 `json:"images"`
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
		fmt.Println(handle.Name(), mime.Value)
		switch {
		case mime.Type == "image":
			image, _ := NewImageFromHandle(handle)
			album.Images = append(album.Images, *image)
		case handle.Name() == "_meta.json":
			err := json.NewDecoder(handle.Descriptor(0)).Decode(&album)
			if err != nil {
				return nil, errors.Wrap(err, "Could not parse album file")
			}
		case mime.Value == "text/markdown":
			fmt.Println("text/markdown")
			md, err := handle.GetUTF8Content()
			if err != nil {
				return nil, errors.Wrap(err, "Could not get content of md file")
			}
			album.Pages = append(album.Pages, md)
		case mime.Value == "text/diary":
			fmt.Println("DIARY")
			album.parseDiary(handle)
		//case strings.HasSuffix(file.Name, ".dia"):
		//	album.parseDiary(file.AsTextfile())
		//		album.parseMeta()
		//		album.Title = fmt.Sprintf("%s | alben", dir.Name)
		case handle.IsDir():
			fmt.Println("dir:", handle.Name())

			var source *Source
			var subalbum *Album
			var images []Image2
			subhandles, err := storage.ReadDir(filepath.Join(path, handle.Name()))
			if err != nil {
				return nil, errors.Wrap(err, "Could not list album subfolder %s", handle.Name())
			}
			for _, subhandle := range subhandles {
				if subhandle.Name() == "_meta.json" {
					subalbum = &Album{Title: handle.Name()}
					err := json.NewDecoder(subhandle.Descriptor(0)).Decode(subalbum)
					if err != nil {
						return nil, errors.Wrap(err, "Could not parse album file")
					}
					fmt.Println("subalbum", subalbum)
				}
				if subhandle.Name() == "_source.json" {
					fmt.Println("subalbum", subhandle.Name())
					source = &Source{Name: handle.Name()}
					err := json.NewDecoder(subhandle.Descriptor(0)).Decode(source)
					if err != nil {
						return nil, errors.Wrap(err, "Could not parse album file")
					}
				}
				submime := subhandle.GuessMIME()
				if submime.Type == "image" {
					subimage, err := NewImageFromHandle(subhandle)
					if err != nil {
						fmt.Println("subalbum error", err)
					}

					images = append(images, *subimage)
				}
			}
			if subalbum != nil {
				subalbum.Images = images

				if len(subalbum.Images) > 0 {
					subalbum.Image = subalbum.Images[0].URL
				}
				album.Moments = append(album.Moments, *subalbum)

				fmt.Println("subalbum images =", subalbum.Images)
			}
			if source != nil {
				source.Images = images
				album.Sources = append(album.Sources, *source)
			}
		}

		//path = filepath.Join(file.Path, handle.Name())

		//folder.Entries = append(folder.Entries, entry)

	}
	fmt.Println("album =", &album.Moments)

	return &album, nil
}

func (a *Album) parseMeta() error {

	//content, err := ioutil.ReadFile(filepath.Join(a.location, "output.json"))

	//json.Unmarshal(content, a)

	return nil
}

func (a *Album) parseDiary(handle Handle) (*Diary, error) {
	fmt.Println("parseDiary", handle.Name())
	diary, err := NewDiary(handle)
	if err != nil {
		fmt.Println("parse diary error", err)
		return nil, err
	}
	a.Diaries = append(a.Diaries, *diary)

	return diary, err
}

// DIARY
type Paragraph struct {
	Content string
}
type Diary struct {
	From       time.Time   `json:"mtime"`
	IndexImage string      `json:"indexImage"`
	Title      string      `json:"title"`
	Content    []Paragraph `json:"content"`
	Images     []Image2    `json:"images"`
}

func NewDiary(handle Handle) (*Diary, error) {

	input := bufio.NewScanner(handle.Descriptor(0))
	diary := &Diary{Title: handle.Name()}
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
			dir := filepath.Dir(handle.URL())
			img := filepath.Join(dir, strings.TrimSpace(line[2:]))
			file, err := handle.Storage().GetHandle(img)
			//			root, e := filepath.Rel(file.location, file.Path)

			image, err := NewImageFromHandle(file) //file.AsImage()
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
