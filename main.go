package main

import (
	"drive/storage"
	"fmt"
	"net/http"
	"os"

	"github.com/namsral/flag"
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
