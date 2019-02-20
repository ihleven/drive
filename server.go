package main

import (
	"drive/gosrc/config"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func CreateServer() {
	fmt.Println(config.Root)
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

	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", logout)
	r.PathPrefix("/serve").HandlerFunc(Serve)
	r.PathPrefix("/").HandlerFunc(PathHandler)

	//r.PathPrefix("/").HandlerFunc(Index)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:29166",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
