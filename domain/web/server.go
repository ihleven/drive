package web

import (
	"drive/config"
	"drive/domain/storage"
	"log"
	"net/http"
	"path"
	"strings"
	"time"
)

func CreateServer() {

	mux := http.NewServeMux()
	mux.Handle("/assets/", assetHandler("assets", "_static/assets"))
	mux.Handle("/dist/", assetHandler("dist", "_static/dist"))
	mux.HandleFunc("/login", Login)
	mux.HandleFunc("/logout", Logout)
	mux.HandleFunc("/serve/home/", Serve(storage.Get("home")))
	mux.HandleFunc("/serve/", Serve(storage.Get("public")))
	mux.HandleFunc("/public/", DispatchStorage(storage.Get("public")))
	mux.HandleFunc("/home/", DispatchStorage(storage.Get("home")))
	mux.HandleFunc("/", IndexHandler)

	srv := &http.Server{
		Handler: mux,
		Addr:    config.Address.String(),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
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
