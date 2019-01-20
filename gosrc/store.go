package main

import (
"fmt"
	"drive/gosrc/config"
)

// Our store will have two methods, to add a new bird,
// and to get all existing birds
// Each method returns an error, in case something goes wrong
type Store interface {
	GetAccommodations() ([]*Accomodation, error)
}

// The store variable is a package level variable that will be available for
// use throughout our application code
var store Store

/*
We will need to call the InitStore method to initialize the store. This will
typically be done at the beginning of our application (in this case, when the server starts up)
This can also be used to set up the store as a mock, which we will be observing
later on
*/
func InitStore(conf config.DatabaseConfiguration) {
	storeConctructor, ok := StoreConstructors[conf.Driver]
	if  !ok {
		fmt.Println("INITStore error:", conf.Driver)
	}
	store, _ = storeConctructor(conf)
	fmt.Println("store:", store)
}

var StoreConstructors map[string]func (config.DatabaseConfiguration) (*mssqlStore, error)

func RegisterStore(key string, constructor func (config.DatabaseConfiguration) (*mssqlStore, error)) {
	if StoreConstructors == nil {
		StoreConstructors =  make(map[string]func (config.DatabaseConfiguration) (*mssqlStore, error))
	}
	StoreConstructors[key] = constructor

}