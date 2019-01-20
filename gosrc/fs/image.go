package fs

import (
	"bufio"
	"drive/gosrc/views"
	"encoding/json"
	"fmt"
	"html/template"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
	"github.com/rwcarlsen/goexif/exif"
)

type Image struct {
	*File
	ColorModel    color.Model
	Width, Height int
	Format        string
	Title         string
	Caption       string // a “caption” is more like a title, while the “cutline” first describes what is happening in the picture, and then explains the significance of the event depicted.
	Cutline       string // the “cutline” is text below a picture, explaining what the reader is looking at

	// https://web.ku.edu/~edit/captions.html
	// https://jerz.setonhill.edu/blog/2014/10/09/writing-a-cutline-three-examples/

	// Caption als allgemeingültige "standalone" Bildunterschrift und Cutline als Verbindung zum Album (ausgewählte Bilder in Reihe?)
}

func NewImage(root, path string) (*Image, error) {
	return nil, nil
}

func (f *File) AsImage() (*Image, error) {

	reader, error := os.Open(f.location)
	if error != nil {
		return nil, error
	}
	defer reader.Close()

	config, format, err := image.DecodeConfig(reader)
	if err != nil {
		log.Fatal(f.Path, config.ColorModel, config.Height, config.Width, err)
		return nil, err
	}
	img := &Image{f, config.ColorModel, config.Width, config.Height, format, "", "", ""}
	img.parseMeta(img.getMetaName())
	fmt.Println("path:", img.Path)
	//img.MakeThumbnail()
	return img, nil

	// exif
	//exif.RegisterParsers(mknote.All...)

	reader.Seek(0, 0)
	x, err := exif.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	camModel, _ := x.Get(exif.Model) // normally, don't ignore errors!
	fmt.Println(camModel.StringVal())

	focal, _ := x.Get(exif.FocalLength)
	numer, denom, _ := focal.Rat2(0) // retrieve first (only) rat. value
	fmt.Printf("%v/%v", numer, denom)

	// Two convenience functions exist for date/time taken and GPS coords:
	tm, _ := x.DateTime()
	fmt.Println("Taken: ", tm)

	lat, long, _ := x.LatLong()
	fmt.Println("lat, long: ", lat, ", ", long)
	j := x.String()
	fmt.Printf("json: %s", j)
	return img, nil

}

func (i Image) getMetaName() string {
	base := strings.TrimSuffix(i.location, filepath.Ext(i.location))
	return fmt.Sprintf("%s.txt", base)
}

func (i *Image) parseMeta(filename string) error {

	fd, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fd.Close()
	s := make([]string, 0)
	r := make(map[string][]string)
	field := "Title"
	input := bufio.NewScanner(fd)
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

	return err
}

func (i *Image) WriteMeta() {
	f, err := os.Create(i.getMetaName())
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

func (i *Image) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPut {
		body, _ := ioutil.ReadAll(r.Body)
		i.update(body)
		json, _ := json.Marshal(i)
		w.Write(json)
	}
	fmt.Println(r.Header.Get("Accept"))
	switch r.Header.Get("Accept") {
	case "application/json":
		views.SerializeJSON(w, i)
	default:
		err := views.Render("image", w, i)
		if err != nil {
			panic(err)
		}
	}
}

func (i *Image) MakeThumbnail() {
	// open "test.jpg"
	file, err := os.Open(i.location)
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

	d, f := filepath.Split(i.location)
	fn := filepath.Join(d, "thumbs", f)
	out, err := os.Create(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)
}
