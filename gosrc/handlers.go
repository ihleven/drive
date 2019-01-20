package main

import (
"drive/gosrc/views"
"net/http"
)

type Accomodation struct {
	AccommodationCode     string `json:"species"`
}

//var accomodations []Accomodation


func GetAccommodationsHandler(w http.ResponseWriter, r *http.Request) {

	accommodations, _ := store.GetAccommodations()

	views.SerializeJSON(w, accommodations)
}