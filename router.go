package main

import (
	"drive/config"
	"drive/fs"
	"drive/models"
	"net/http"
	"path"
	"regexp"
	"strings"
)

type Muxer struct {
	routes map[string]http.Handler
}

func (m *Muxer) register(pattern string, cb http.Handler) {
	m.routes[pattern] = cb

}

var digitsRegexp = regexp.MustCompile(`foo/(?P<second>\d+)`)
var helloRE = regexp.MustCompile(`/hello/(?P<second>\w+)`)

func (p *Muxer) ServeHTTP(w http.ResponseWriter, request *http.Request) {

	route := path.Clean(request.URL.Path)
	//elements := strings.SplitN(route, "/", 4)

	if Filer, error := fs.NewFiler(config.Root, route); error == nil {

		Filer.Render(w, request)

	} else if remainder := strings.TrimPrefix(route, "/drive"); len(remainder) < len(route) {
		models.PathHandler(w, request, remainder)

	} else if m := helloRE.FindStringSubmatch(route); m != nil {
		//sayhelloName(w, request, m[1])
	} else {
		http.NotFound(w, request)

	}
}

var countryCode = map[string]string{
	"deutschland": "de",
	"frankreich":  "fr",
	"germany":     "de"}

func isCountry(token string) (bool, string) {
	code, ok := countryCode[token]
	return ok, code
}
