package models

import (
	"bufio"
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
		format, "Ich bin's", "", "", nil,
	}
	image.parseMeta()
	//img.MakeThumbnail()

	if err = image.GoexifDecode(); err != nil {
		fmt.Println("Error Decoding Exif with Goexif =>", err)
	}

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
		log.Fatal(err)
	}

	camModel, _ := x.Get(exif.Model) // normally, don't ignore errors!
	model, _ := camModel.StringVal()
	e.Model = model

	orientation, _ := x.Get(exif.Orientation)
	o, err := orientation.Int(0)
	e.Orientation = o

	focal, _ := x.Get(exif.FocalLength)
	numer, denom, _ := focal.Rat2(0) // retrieve first (only) rat. value
	fmt.Printf("%v/%v %s\n", numer, denom, focal.String())

	// Two convenience functions exist for date/time taken and GPS coords:
	tm, _ := x.DateTime()
	e.Taken = tm

	lat, long, _ := x.LatLong()
	e.Lat = lat
	e.Lng = long

	//j := x.String()
	//fmt.Printf("json: %s", j)

	return nil
}

func (i Image) getMetaFile() *file.Info {
	base := strings.TrimSuffix(i.File.Path, filepath.Ext(i.File.Path))
	f, _ := file.Open(fmt.Sprintf("%s.txt", base), 0, 0, 0)
	return f
}

func (i *Image) parseMeta() error {
	fmt.Println("parseMeta")
	meta := i.getMetaFile()
	defer meta.Close()

	s := make([]string, 0)
	r := make(map[string][]string)
	field := "Title"
	fmt.Println("title")
	input := bufio.NewScanner(meta.Descriptor)
	for input.Scan() {
		line := strings.TrimSpace(input.Text())
		if line == "" {
			continue
		} else if strings.Trim(line, "=") == "" {
			field = "Caption"
			continue
		} else if strings.Trim(line, "-") == "" {
			if field == "Caption" {
				field = "Cutline"
			} else {
				field = "NExt"
			}
			continue
		}
		s = append(s, line)
		r[field] = append(r[field], line)
		fmt.Println(field, line)
	}
	i.Title = r["Title"][0]
	i.Caption = r["Caption"][0]
	i.Cutline = strings.Join(r["Cutline"], "\n")
	fmt.Println(i.Title, i.Caption, i.Cutline)
	return nil
}

func (i *Image) WriteMeta() {
	f, err := os.Create(i.getMetaFile().Name())
	if err != nil {
		panic(err)
	}
	tmpl, err := template.New("txt").Parse("{{.Title}}\n=====\n{{.Caption}}\n-----\n{{.Cutline}}\n------\n")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(f, i)
	if err != nil {
		panic(err)
	}
	f.Close()
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

	i.WriteMeta()
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
