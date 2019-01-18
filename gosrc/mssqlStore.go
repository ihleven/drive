package main

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
"fmt"
"log"
	_ "github.com/lib/pq"
	"drive/gosrc/config"
)

func init() {
	RegisterStore("sqlserver", NewMSSQLStore)
}

type mssqlStore struct {
	db *sql.DB
}


func NewMSSQLStore(config config.DatabaseConfiguration) (*mssqlStore, error) {

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;encrypt=disable", config.Host, config.User, config.Password, config.Port, config.Name)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return &mssqlStore{db}, nil
}


func (store *mssqlStore) GetAccommodations() ([]*Accomodation, error) {

	rows, err := store.db.Query("SELECT AccommodationCode FROM [webcc-local].[rx].[Accommodations]")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	birds := []*Accomodation{}
	for rows.Next() {
		bird := &Accomodation{}
		if err := rows.Scan(&bird.AccommodationCode); err != nil {
			return nil, err
		}
		birds = append(birds, bird)
	}
	return birds, nil
}