package models

import (
	"drive/domain"
	"drive/domain/storage"
	"drive/file"
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
	File          *file.File
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
	Exif *Exif
}
type Exif struct {
	Orientation int
	Taken       time.Time
	Lat,
	Lng float64
	Model string
}

func NewImage(file *file.File) (*Image, error) {
	config, format, err := image.DecodeConfig(file.Descriptor)
	if err != nil {
		log.Fatal(file.Path, config.ColorModel, config.Height, config.Width, err)
		return nil, err
	}
	image := &Image{file,
		config.ColorModel, config.Width, config.Height, float64(config.Height) / float64(config.Width) * 100,
		format, "", "", "", nil,
	}
	if err = image.parseMeta(); err != nil {
		fmt.Println("Error parsing meta =>", err)
	}
	//img.MakeThumbnail()
	if err = image.GoexifDecode(); err != nil {
		fmt.Println("Error Decoding Exif with Goexif =>", err)
	}
	fmt.Println("asdf3")
	return image, nil

}

func (i *Image) GoexifDecode() error {
	// https://github.com/rwcarlsen/goexif

	// Optionally register camera makenote data parsing - currently Nikon and
	// Canon are supported.
	//exif.RegisterParsers(mknote.All...)

	e := &Exif{}
	i.Exif = e

	i.File.Descriptor.Seek(0, 0)
	x, err := exif.Decode(i.File.Descriptor)
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
	e.Model = model

	orientation, err := x.Get(exif.Orientation)
	if err != nil {
		return err
	}
	o, err := orientation.Int(0)

	if err != nil {
		return err
	}
	e.Orientation = o

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
	e.Taken = tm

	lat, long, err := x.LatLong()
	if err != nil {
		return err
	}
	e.Lat = lat
	e.Lng = long

	//j := x.String()
	//fmt.Printf("json: %s", j)

	return nil
}

func (i Image) getMetaFile() *file.Info {
	base := strings.TrimSuffix(i.File.Path, filepath.Ext(i.File.Path))
	f, err := file.Open(fmt.Sprintf("%s.txt", base), 0, 0, 0)
	if err != nil {
		fmt.Println("error opening meta file", err)
		return nil
	}
	return f
}

func (i *Image) parseMeta() error {

	re := regexp.MustCompile(`(?s)(?P<Title>.*?)=+(?P<Caption>.*?)---+(?P<Cutline>.*?)---+`)

	f, err := storage.DefaultStorage.OpenFile(i.getMetaFilename(), os.O_RDONLY, 0)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	content, err := f.GetContent()
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

func (i Image) getMetaFilename() string {
	base := strings.TrimSuffix(i.File.Path, filepath.Ext(i.File.Path))
	filename := fmt.Sprintf("%s.txt", base)
	return filename
}

func (i *Image) WriteMeta(usr *domain.Account) error {
	filename := i.getMetaFilename()
	fmt.Println("storage.DefaultStorage.OpenFile", filename)

	f, err := storage.DefaultStorage.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	tmpl, err := template.New("txt").Parse("{{.Title}}\n=====\n{{.Caption}}\n-----\n{{.Cutline}}\n------\n")
	if err != nil {
		return err
	}
	err = tmpl.Execute(f.File, i)
	if err != nil {
		return err
	}
	f.Close()
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

func (i *Image) MakeThumbnail() {
	// open "test.jpg"
	file, err := os.Open(i.File.Path)
	if err != nil {
		log.Fatal(err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(100, 0, img, resize.Lanczos3)

	d, f := filepath.Split(i.File.Path)
	fn := filepath.Join(d, "thumbs", f)
	out, err := os.Create(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)
}
