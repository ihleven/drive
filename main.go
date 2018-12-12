package main

import (
	"context"
	"database/sql"
	"drive/config"
	"fmt"
	"log"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb"
)

//const serverUA = "Alexis/0.2"
const fs_maxbufsize = 4096 // 4096 bits = default page size on OSX

var dirlisting_tpl = "templates/directory.html"

// Manages directory listings
type dirlisting struct {
	Name           string
	Children_dir   []string
	Children_files []string
	ServerUA       string
}

func main() {
	//dbf()
	config.ParseFlags()

	mux := &Muxer{}
	//mux.register("/blah", http.HandlerFunc(sayhelloName))
	http.ListenAndServe(config.Address.String(), mux)
	http.ListenAndServe(config.Address.String(), http.HandlerFunc(pathRequestHandler))
}

/* Go is the first programming language with a templating engine embeddeed
 * but with no min function. */
func min(x int64, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

var db *sql.DB

func dbf() {

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;encrypt=disable",
		config.Dbconn.Host, config.Dbconn.User, config.Dbconn.Password, config.Dbconn.Port, config.Dbconn.Name)

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
