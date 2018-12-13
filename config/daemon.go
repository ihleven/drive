package config

import (
	"fmt"
	"os"

	"github.com/namsral/flag"
)

//const serverUA = "Alexis/0.2"
const fs_maxbufsize = 4096 // 4096 bits = default page size on OSX

var (
	verbose bool
	Address address
	Root    string
	cwd     string
)

func ParseFlags() {

	//cwd, _ := os.Getwd()

	flag.StringVar(&Address.host, "host", "localhost", "Host")
	flag.IntVar(&Address.port, "port", 3000, "Port")
	flag.StringVar(&Root, "root", getCwd(), "Root folder")
	flag.BoolVar(&verbose, "verbose", false, "help message")
	flag.Parse()
}

type address struct {
	host string
	port int
}

func (a *address) String() string {
	return fmt.Sprintf("%s:%d", Address.host, Address.port)
}

func getCwd() string {
	return "/Users/mi/tmp/"
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error while getting current directory.")
		return ""
	}
	return cwd
}
