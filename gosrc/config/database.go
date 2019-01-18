package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"

)

type DatabaseConfiguration struct {
	Driver string
	Name     string
	User     string
	Password string
	Host     string
	Port     int
}

var databases = map[string]DatabaseConfiguration{
	"mssql": DatabaseConfiguration{Driver: "sqlserver", Name: "webcc-local", Host: "127.0.0.1", Port: 1433, User: "webrx", Password: "0Q2u09KnbnawxEDEtwox"}}
 

var db *sql.DB

func GetDatabaseConfiguration(key string) DatabaseConfiguration {

	conf, ok := databases[key]
	if  !ok {
		panic(fmt.Sprintf("configuration not found: %s", key))
	}
	return conf
}

func asdf() {
	//alt
	ctx := context.Background()
	err := db.PingContext(ctx)
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
