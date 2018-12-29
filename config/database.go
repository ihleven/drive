package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

type DatabaseConnection struct {
	Name     string
	Host     string
	Port     int
	User     string
	Password string
}

var Dbconn = &DatabaseConnection{Name: "webcc-local", Host: "127.0.0.1", Port: 1433, User: "webrx", Password: "0Q2u09KnbnawxEDEtwox"}

var db *sql.DB

func dbf() {

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;encrypt=disable",
		Dbconn.Host, Dbconn.User, Dbconn.Password, Dbconn.Port, Dbconn.Name)

	var err error

	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Connected!\n")

	//ctx = context.Background()

	// Check if database is alive.
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	tsql := fmt.Sprintf("SELECT AccommodationCode FROM [webcc-local].[rx].[Accommodations];")

	// Execute query
	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer rows.Close()

	var count int

	// Iterate through the result set.
	for rows.Next() {
		var code string
		//var id int

		// Get values from row.
		err := rows.Scan(&code)
		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Printf("ID: %s\n", code)
		count++
	}

	log.Fatal(err.Error())
}
