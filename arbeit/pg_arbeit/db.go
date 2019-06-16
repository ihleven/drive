package pg_arbeit

import (
	"database/sql"
	"drive/config"
	"drive/errors"
	"fmt"
)

type database struct {
	DB *sql.DB
}

func (db database) Close() {
	db.DB.Close()
}

var repo database

// GetDatabaseHandle creates and verifies a database handle, returning it to the caller
func GetDatabaseHandle(conf config.DatabaseConfiguration) (*database, error) {

	dbinfo := fmt.Sprintf("host=localhost port=5432 user=%s password=%s dbname=%s sslmode=disable", conf.User, conf.Password, conf.Name)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return nil, errors.Augment(err, errors.BadCredentials, "Could not get database handle '%s'", dbinfo)
	}
	// kann hier nicht geschlossen werde, weil das handle zurueckgegeben wird.
	//defer db.Close()

	err = db.Ping()
	//Ping verifies a connection to the database is still alive, establishing a connection if necessary.
	if err != nil {
		return nil, errors.Wrap(err, "Could not verify database connection.")
	}
	repo = database{db}
	return &repo, nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
