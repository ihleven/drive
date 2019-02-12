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
	r.PathPrefix("/serve/").Handler(
		http.StripPrefix("/serve/", http.FileServer(http.Dir(config.Root))),
	)
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("./vue/node_modules/semantic-ui-css"))),
	)
	r.PathPrefix("/dist/").Handler(
		http.StripPrefix("/dist/", http.FileServer(http.Dir("./vue/dist"))),
	)
	r.PathPrefix("/assets/").Handler(
		http.StripPrefix("/assets/", http.FileServer(http.Dir("./vue/dist/assets"))),
	)

	r.HandleFunc("/secret", secret)
	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", logout)

	r.PathPrefix("/").HandlerFunc(Index)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:3001",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
