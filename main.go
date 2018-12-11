package main

import (
	"drive/storage"
	"fmt"
	"net/http"
	"os"
	_ "github.com/denisenkom/go-mssqldb"
    "database/sql"
    "log"
	"github.com/namsral/flag"
	"context"
	"drive/config"
)

var root_folder *string // TODO: Find a way to be cleaner !
var uses_gzip *bool

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
	dbf()
	// Get current working directory to get the file from it
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error while getting current directory.")
		return
	}
	fmt.Println(cwd)
	var port int
	flag.IntVar(&port, "port", 3001, "Port number")
	//flag.Parse()

	// Command line parsing
	bind := flag.String("bind", ":3000", "Bind address")
	ro := flag.String("root", "/Users/mi/tmp/", "Root folder")
	storage.Root_folder = *ro
	uses_gzip = flag.Bool("gzip", true, "Enables gzip/zlib compression")

	flag.Parse()

	//http.Handle("/", )

	fmt.Printf("Sharing %s on %d ...\n", storage.Root_folder, port)
	mux := &Muxer{}
	//mux.register("/blah", http.HandlerFunc(sayhelloName))
	http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	http.ListenAndServe((*bind), http.HandlerFunc(pathRequestHandler))
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


//jdbc:sqlserver://127.0.0.1:1433;databaseName=webcc-local;user=webrx;password=0Q2u09KnbnawxEDEtwox
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