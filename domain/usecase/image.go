package usecase

import (
	"drive/domain"
	"encoding/json"
	"errors"
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
	File          *domain.File
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
	metaFile *domain.File
}
type Exif struct {
	Orientation int
	Taken       time.Time
	Lat,
	Lng float64
	Model string
}

func NewImage(file *domain.File, usr *domain.Account) (*Image, error) {

	fd := file.Descriptor()
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

	//img.MakeThumbnail()
	if err = i.GoexifDecode(fd); err != nil {
		fmt.Println("Error Decoding Exif with Goexif =>", err)
	}
	fmt.Println("asdf3", i.Exif)
	return i, nil
}

func (i *Image) getMetaFilename() string {
	base := strings.TrimSuffix(i.File.Path, filepath.Ext(i.File.Path))
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
	i.Exif.Orientation = o

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
	i.Exif.Taken = tm

	lat, long, err := x.LatLong()
	if err != nil {
		return err
	}
	i.Exif.Lat = lat
	i.Exif.Lng = long

	//j := x.String()
	//fmt.Printf("json: %s", j)

	return nil
}

func (i *Image) WriteMeta(usr *domain.Account) error {

	if !i.metaFile.Permissions.Write {
		return errors.New(fmt.Sprintf("Missing write permission for %s", i.metaFile.Name))
	}
	fd := i.metaFile.Descriptor()
	fd.Close()

	tmpl, err := template.New("txt").Parse("{{.Title}}\n=====\n{{.Caption}}\n-----\n{{.Cutline}}\n------\n")
	if err != nil {
		return err
	}
	err = tmpl.Execute(fd, i)
	if err != nil {
		return err
	}
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
