package config

import (
	"fmt"
)

type settings struct {
	Address        address
	Root           string
	cwd            string
	configFilePath string
	Storages       map[string]storage
}
type address struct {
	host string
	port int
}

var Settings settings

func (a *address) String() string {
	return fmt.Sprintf("%s:%d", a.host, a.port)
}
