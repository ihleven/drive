package drive

import (
	"drive/domain"
	"drive/errors"
	"encoding/json"
	"fmt"
	"html/template"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/nfnt/resize"
	"github.com/rwcarlsen/goexif/exif"
)

type Image struct {
	File          *File
	ColorModel    color.Model
	Width, Height int
	Ratio         float64
	Format        string
	Title         string
	Caption       string // a “caption” is more like a title, while the “cutline” first describes what is happening in the picture, and then explains the significance of the event depicted.
	Cutline       string // the “cutline” is text below a picture, explaining what the reader is looking at

	// https://web.ku.edu/~edit/captions.html
	// https://jerz.setonhill.edu/blog/2014/10/09/writing-a-cutline-three-examples/

	// Caption als allgemeingültige "standalone" Bildunterschrift und Cutline als Verbindung zum Album (ausgewählte Bilder in Reihe?)
	Exif     *Exif
	metaFile *File
}

type Image2 struct {
	Handle     Handle      `json:"_"`
	ColorModel color.Model `json:"_"`
	Width      int         `json:"w"`
	Height     int         `json:"h"`
	Ratio      float64     `json:"_"`
	Format     string      `json:"_"`
	Title      string      `json:"title,omitempty"`
	Caption    string      `json:"_"` // a “caption” is more like a title, while the “cutline” first describes what is happening in the picture, and then explains the significance of the event depicted.
	Cutline    string      `json:"_"` // the “cutline” is text below a picture, explaining what the reader is looking at

	// https://web.ku.edu/~edit/captions.html
	// https://jerz.setonhill.edu/blog/2014/10/09/writing-a-cutline-three-examples/

	// Caption als allgemeingültige "standalone" Bildunterschrift und Cutline als Verbindung zum Album (ausgewählte Bilder in Reihe?)
	Exif     *Exif  `json:"_"`
	MetaFile *File  `json:"_"`
	Src      string `json:"-"` // `json:"src"`
	Name     string `json:"name"`
	Source   string `json:"source"`
}
type Exif struct {
	Orientation *int
	Taken       *time.Time
	Lat,
	Lng *float64
	Model string
}

func NewImageFromHandle(handle Handle) (*Image2, error) {

	fd := handle.Descriptor(0)
	defer fd.Close()

	config, format, err := image.DecodeConfig(fd)
	if err != nil {
		return nil, errors.Wrap(err, "Error decoding image")
	}

	i := &Image2{
		Handle:     handle,
		ColorModel: config.ColorModel,
		Width:      config.Width,
		Height:     config.Height,
		Ratio:      float64(config.Height) / float64(config.Width) * 100,
		Format:     format,
		Src:        handle.StoragePath(),
		Name:       handle.Name(),
	}
	return i, nil
}

func NewImageFromHandle2(handle Handle, prefix string) (*Image2, error) {

	fd := handle.Descriptor(0)
	defer fd.Close()

	config, format, err := image.DecodeConfig(fd)
	if err != nil {
		return nil, errors.Wrap(err, "Error decoding image")
	}

	i := &Image2{
		Handle:     handle,
		ColorModel: config.ColorModel,
		Width:      config.Width,
		Height:     config.Height,
		Ratio:      float64(config.Height) / float64(config.Width) * 100,
		Format:     format,
		//Src:        handle.ServeURL(),
		Name:   handle.Name(),
		Source: prefix,
	}
	return i, nil
}

func NewImage(file *File, usr *domain.Account) (*Image, error) {

	fd := file.Descriptor(0)
	defer fd.Close()
	config, format, err := image.DecodeConfig(fd)
	if err != nil {
		log.Fatal("NewImage", err)
		return nil, err
	}

	i := &Image{
		File:       file,
		ColorModel: config.ColorModel,
		Width:      config.Width,
		Height:     config.Height,
		Ratio:      float64(config.Height) / float64(config.Width) * 100,
		Format:     format,
	}

	metafile, err := GetFile(file.Storage(), i.getMetaFilename(), usr)
	if err == nil {
		i.metaFile = metafile

		if err = i.parseMeta(); err != nil {
			fmt.Println("Error parsing meta =>", err)
		}
	}
	fmt.Println("NewImage", i.getMetaFilename(), fd.Name())

	//img.MakeThumbnail()
	if err = i.GoexifDecode(fd); err != nil {
		fmt.Println("Error Decoding Exif with Goexif =>", err)
	}
	fmt.Println("asdf3", i.Exif, err)
	return i, nil
}

func (i *Image) getMetaFilename() string {
	base := strings.TrimSuffix(i.File.StoragePath(), filepath.Ext(i.File.StoragePath()))
	filename := fmt.Sprintf("%s.txt", base)
	return filename
}

func (i *Image) parseMeta() error {

	re := regexp.MustCompile(`(?s)(?P<Title>.*?)=+(?P<Caption>.*?)---+(?P<Cutline>.*?)---+`)

	content, err := i.metaFile.GetContent()
	if err != nil {
		return err
	}

	match := re.FindSubmatch(content)
	paramsMap := make(map[string]string)

	for i, name := range re.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = strings.TrimSpace(string(match[i]))
		}
	}
	if title, ok := paramsMap["Title"]; ok {
		i.Title = title
	}
	if caption, ok := paramsMap["Caption"]; ok {
		i.Caption = caption
	}
	if cutline, ok := paramsMap["Cutline"]; ok {
		i.Cutline = cutline
	}
	return nil
}

func (i *Image) GoexifDecode(fd *os.File) error {
	// https://github.com/rwcarlsen/goexif

	// Optionally register camera makenote data parsing - currently Nikon and
	// Canon are supported.
	//exif.RegisterParsers(mknote.All...)

	i.Exif = &Exif{}

	fd.Seek(0, 0)
	x, err := exif.Decode(fd)
	if err != nil {
		return err
	}

	camModel, err := x.Get(exif.Model) // normally, don't ignore errors!
	if err != nil {
		return err
	}
	model, err := camModel.StringVal()
	if err != nil {
		return err
	}
	i.Exif.Model = model

	orientation, err := x.Get(exif.Orientation)
	if err != nil {
		return err
	}
	o, err := orientation.Int(0)

	if err != nil {
		return err
	}
	i.Exif.Orientation = &o

	focal, _ := x.Get(exif.FocalLength)
	numer, denom, err := focal.Rat2(0) // retrieve first (only) rat. value
	if err != nil {
		return err
	}
	fmt.Printf("%v/%v %s\n", numer, denom, focal.String())

	// Two convenience functions exist for date/time taken and GPS coords:
	tm, err := x.DateTime()
	if err != nil {
		return err
	}
	i.Exif.Taken = &tm

	lat, long, err := x.LatLong()
	if err != nil {
		return err
	}
	i.Exif.Lat = &lat
	i.Exif.Lng = &long

	//j := x.String()
	//fmt.Printf("json: %s", j)

	return nil
}

func (i *Image) WriteMeta(usr *domain.Account) error {
	fmt.Println("writemata1", i.metaFile)

	if i.metaFile == nil {
		filename := i.getMetaFilename()
		file, err := CreateFile(i.File.Storage(), filename, usr)
		if err != nil {
			return errors.Errorf("Could not create meta file: %s", filename)
		}
		i.metaFile = file
	}

	if !i.metaFile.Permissions.Write {
		return errors.Errorf("Missing write permission for %s", i.metaFile.Name)
	}
	fd := i.metaFile.Descriptor(os.O_CREATE | os.O_WRONLY | os.O_TRUNC)
	defer fd.Close()
	fmt.Println("writemeta", fd.Name())

	tmpl, err := template.New("txt").Parse("{{.Title}}\n=====\n{{.Caption}}\n-----\n{{.Cutline}}\n------\n")
	if err != nil {
		return err
	}

	err = tmpl.Execute(fd, i)
	if err != nil {
		return err
	}
	fmt.Println("writemeta", fd.Name())
	return nil
}

func (i *Image) update(requestBody []byte) {

	var m map[string]interface{}
	_ = json.Unmarshal(requestBody, &m)

	if title, ok := m["Title"]; ok {
		i.Title = title.(string)
		fmt.Printf(" - update Title => '%s'\n", title)
	}
	if caption, ok := m["Caption"]; ok {
		i.Caption = caption.(string)
		fmt.Printf(" - update Caption => '%s'\n", caption)
	}
	if cutline, ok := m["Cutline"]; ok {
		i.Cutline = cutline.(string)
		fmt.Printf(" - update Cutline => '%s'\n", cutline)
	}

	i.WriteMeta(nil)
}

func MakeThumbs(handle Handle) error {

	thumbs := filepath.Join(handle.Location(), "thumbs")
	if err := os.RemoveAll(thumbs); err != nil {
		return errors.Wrap(err, "Could not delete thumbs folder")
	}
	if err := os.Mkdir(thumbs, 0755); err != nil {
		return errors.Wrap(err, "Could not create thumbs folder")
	}
	for _, format := range thumbFormats {
		if err := os.Mkdir(filepath.Join(thumbs, format.Name), 0755); err != nil {
			return errors.Wrap(err, "Could not create thumbs folder for format %v", format)
		}
	}
	handles, err := handle.Storage().ReadDir(handle.StoragePath())
	if err != nil {
		return errors.Wrap(err, "Could not read thumbs folder")
	}

	for _, entry := range handles {

		if mime := entry.GuessMIME(); mime.Value == "image/jpeg" {

			err = MakeThumbnail(entry, thumbs)
			if err != nil {
				return errors.Wrap(err, "Could not create thumbnail")
			}
		}

	}
	return nil
}

type thumbformat struct {
	Name string
	X, Y uint
}

var thumbFormats = []thumbformat{{"x100", 0, 100}, {"96", 96, 96}} //map[string]thumbformat

func MakeThumbnail(orig Handle, thumbnailFolder string) error {

	fd := orig.Descriptor(0)
	defer fd.Close()
	if fd == nil {
		return errors.Errorf("Could not open fd")
	}
	// decode jpeg into image.Image
	img, err := jpeg.Decode(fd)
	if err != nil {
		return errors.Wrap(err, "error jpeg decoding %s", orig.Name())
	}

	for _, format := range thumbFormats {
		fmt.Println(format)
		// resize to width 100 using Lanczos resampling
		// and preserve aspect ratio
		thumb := resize.Resize(format.X, format.Y, img, resize.Lanczos3)

		location := filepath.Join(thumbnailFolder, format.Name, orig.Name())
		fmt.Println(format, location)
		dest, err := os.OpenFile(location, os.O_RDWR|os.O_CREATE, 0755)
		defer dest.Close()
		if err != nil {
			return errors.Wrap(err, "Failed to create file: %s", location)
		}

		// write new image to file
		err = jpeg.Encode(dest, thumb, nil)
		if err != nil {
			return errors.Wrap(err, "Could not jpeg.Encode")
		}
	}
	return nil
}
