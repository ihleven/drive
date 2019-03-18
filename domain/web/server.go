package web

import (
	"drive/config"
	"drive/domain/storage"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func CreateServer() {

	r := mux.NewRouter()
	// This will serve files under http://localhost:8000/static/<filename>

	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))),
	)
	r.PathPrefix("/dist/").Handler(
		http.StripPrefix("/dist/", http.FileServer(http.Dir("./vue/dist"))),
	)
	r.PathPrefix("/assets/").Handler(
		http.StripPrefix("/assets/", http.FileServer(http.Dir("./vue/dist/assets"))),
	)

	//r.HandleFunc("/login", handler.Login)
	//r.HandleFunc("/logout", handler.Logout)
	//r.PathPrefix("/serve").HandlerFunc(handler.Raw)
	//r.PathPrefix("/").Handler(pathRouter) //PathHandler)

	//r.PathPrefix("/").HandlerFunc(Index)

	mux := http.NewServeMux()
	mux.Handle("/static/", assetHandler("static", "static"))
	mux.Handle("/dist/", assetHandler("dist", "vue/dist"))
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
