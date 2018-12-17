package main

import (
	"context"
	"database/sql"
	"drive/config"
	"drive/templates"
	"drive/views"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	//dbf()
	config.ParseFlags()
	templates.Init()

	http.Handle("/serve/", http.StripPrefix("/serve/", http.FileServer(http.Dir(config.Root))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/hello/", http.StripPrefix("/hello/", http.HandlerFunc(sayhelloName)))
	http.HandleFunc("/hallo/", sayhelloName)
	http.HandleFunc("/", views.PathHandler)

	mux := &Muxer{}
	http.Handle("/drive", http.StripPrefix("/drive", mux))
	//router := Router{}
	//
	http.ListenAndServe(config.Address.String(), nil)
	//http.ListenAndServe(config.Address.String(), http.HandlerFunc(pathRequestHandler))
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Hello %s", r.URL.Path)

}

/* Go is the first programming language with a templating engine embeddeed
 * but with no min function. */
func min(x int64, y int64) int64 {
	if x < y {
		return x
	}
	return y
}
func parseCSV(data string) []string {
	splitted := strings.SplitN(data, ",", -1)

	data_tmp := make([]string, len(splitted))

	for i, val := range splitted {
		data_tmp[i] = strings.TrimSpace(val)
	}

	return data_tmp
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
