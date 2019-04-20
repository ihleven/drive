package web

import (
	"drive/config"
	"log"
	"net/http"
	"path"
	"strings"
	"time"
)

var mux *http.ServeMux

func init() {

	mux = http.NewServeMux()
}

func CreateServer() {

	mux.Handle("/assets/", assetHandler("assets", "_static/assets"))
	mux.Handle("/dist/", assetHandler("dist", "_static/dist"))
	mux.HandleFunc("/login", Login)
	mux.HandleFunc("/logout", Logout)
	mux.HandleFunc("/", defaultHandler)

	srv := &http.Server{
		Handler: mux,
		Addr:    config.Address.String(),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func RegisterFunc(pattern string, handlerfunc func(w http.ResponseWriter, r *http.Request)) {
	mux.HandleFunc(pattern, handlerfunc)
}

func RegisterHandler(pattern string, handler http.Handler) {
	mux.Handle(pattern, handler)
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
