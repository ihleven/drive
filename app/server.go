package app

import (
	"drive/config"
	"drive/handler"
	"log"
	"net/http"
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

	r.HandleFunc("/login", handler.Login)
	r.HandleFunc("/logout", handler.Logout)
	r.PathPrefix("/serve").HandlerFunc(handler.Serve)
	r.PathPrefix("/").Handler(pathRouter) //PathHandler)

	//r.PathPrefix("/").HandlerFunc(Index)

	srv := &http.Server{
		Handler: r,
		Addr:    config.Address.String(),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
