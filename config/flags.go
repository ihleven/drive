package config

import (
	"fmt"
	"os"

	"github.com/namsral/flag"
)

func ParseFlags() {

	//cwd, _ := os.Getwd()

	flag.StringVar(&Address.host, "host", "localhost", "Host")
	flag.IntVar(&Address.port, "port", 3000, "Port")
	flag.StringVar(&Root, "root", getCwd(), "Root folder")
	flag.BoolVar(&verbose, "verbose", false, "help message")
	flag.Parse()
	fmt.Println(Root)
}

func getCwd() string {

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error while getting current directory.")
		return ""
	}
	return cwd
}
