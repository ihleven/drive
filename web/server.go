package web

import (
	"log"
	"net/http"
	"path"
	"strings"
	"time"
)

type server struct {
	address string
	mux     *http.ServeMux
}

func NewServer(address string) *server {

	mux := http.NewServeMux()

	mux.Handle("/assets/", assetHandler("assets", "_static/assets"))
	mux.Handle("/dist/", assetHandler("dist", "_static/dist"))
	mux.HandleFunc("/login", Login)
	mux.HandleFunc("/logout", Logout)
	mux.HandleFunc("/", defaultHandler)

	return &server{address: address, mux: mux}
}

func (s server) Run() {

	httpServer := &http.Server{
		Handler: s.mux,
		Addr:    s.address,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(httpServer.ListenAndServe())
}

func (s server) RegisterHandlerFunc(pattern string, handlerfunc func(w http.ResponseWriter, r *http.Request)) {
	s.mux.HandleFunc(pattern, handlerfunc)
}

func (s server) RegisterHandler(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

func logit(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Before")
		h.ServeHTTP(w, r) // call original
		log.Println("After")
	})
}

// ShiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
