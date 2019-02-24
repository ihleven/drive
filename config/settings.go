package config

import (
	"fmt"
)

//const serverUA = "Alexis/0.2"
const Fs_maxbufsize = 4096 // 4096 bits = default page size on OSX

var (
	verbose bool
	Address address
	Root    string
	cwd     string
)

type address struct {
	host string
	port int
}

func (a *address) String() string {
	return fmt.Sprintf("%s:%d", Address.host, Address.port)
}
