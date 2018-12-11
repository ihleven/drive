package daemon

import (
	"drive/storage"
	"fmt"
	"net/http"
	"path"
	"regexp"
	"strings"
	"sync"
)





type Router struct { // reuse ServeMux type of net/http
	mu    sync.RWMutex
	m     map[string]http.Handler
}

func (r *Router) Register(pattern string, handler http.Handler) {
	// Geklaut von net/http.Handle => Handle registers the handler for the given pattern.
	// If a handler already exists for pattern, Handle panics.
	r.mu.Lock()
	defer r.mu.Unlock()

	if pattern == "" || handler == nil {
		panic("http: invalid pattern OR nil handler")
	}
	if _, exist := r.m[pattern]; exist {
		panic("http: multiple registrations for " + pattern)
	}

	if r.m == nil {
		r.m = make(map[string]http.Handler)
	}
	r.m[pattern] = handler
}
// HandleFunc registers the handler function for the given pattern.
func (r *Router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	if handler == nil {
		panic("http: nil handler")
	}
	r.Register(pattern, http.HandlerFunc(handler))
}


var digitsRegexp = regexp.MustCompile(`foo/(?P<second>\d+)`)
var helloRE = regexp.MustCompile(`/hello/(?P<second>\w+)`)

func (r *Router) ServeHTTP(w http.ResponseWriter, request *http.Request) {

	route := path.Clean(request.URL.Path)
	elements := strings.SplitN(route, "/", 4)

	if valid := storage.ValidPath(route); valid {

		fmt.Printf(" - serving: %s", route)
		p := fmt.Sprintf("/Users/mi/%s", route)
		//fs := http.FileServer(http.Dir("/Users/mi/"))
		http.ServeFile(w, request, p)
		//fs.ServeHTTP(w, request)

	} else if remainder := strings.TrimPrefix(route, "/drive"); len(remainder) < len(route) {
		storage.PathHandler(w, request, remainder)

	} else if albumname := strings.TrimPrefix(route, "/alben"); len(albumname) < len(route) {
		storage.HandleAlbum(w, request, albumname)
	} else
	//hello/
	if m := helloRE.FindStringSubmatch(route); m != nil {
		sayhelloName(w, request, m[1])
	} else if code, ok := countryCode[elements[1]]; ok {
		search(w, request, code)
		return
	} else if request.URL.Path == "/hello" {
		sayhelloName(w, request, "d")
		return
	} else if baseFile := storage.GetBaseFile(route); baseFile != nil {
		if baseFile.Info.IsDir() {
			//return views.Directory(w, request, baseFile)
			dir, _ := storage.NewDirectory(baseFile.Info, baseFile.Path)
			dir.RenderHTML(w, request)
		}
		http.ServeFile(w, request, baseFile.Path)
		return
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

func sayhelloName(w http.ResponseWriter, r *http.Request, name string) {
	fmt.Fprintf(w, "Hello %s", name)
}
func search(w http.ResponseWriter, r *http.Request, country string) {
	fmt.Fprintf(w, "Region %s!", country)
}
