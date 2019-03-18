package config

import (
	"fmt"
	"time"
)

//const serverUA = "Alexis/0.2"
const Fs_maxbufsize = 4096 // 4096 bits = default page size on OSX

var (
	verbose bool
	Address address
	Root    string
	cwd     string
)
var conf Config

type Config struct {
	Public     string
	Homes      []string
	Pi         float64
	Perfection []int
	DOB        time.Time // requires `import time`
}

type address struct {
	host string
	port int
}

func (a *address) String() string {
	return fmt.Sprintf("%s:%d", Address.host, Address.port)
}
