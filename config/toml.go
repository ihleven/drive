package config

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type tomlConfig struct {
	Title    string `toml:"title"`
	Server   server
	DB       database `toml:"database"`
	Storages map[string]storage
	Clients  clients
}

type server struct {
	Host string
	Port int `toml:"port"`
}

type database struct {
	Server  string
	Ports   []int
	ConnMax int `toml:"connection_max"`
	Enabled bool
}

type storage struct {
	Root      string
	BaseURL   string
	ServeURL  string
	AlbumPath string
	Group     uint32
	Mode      os.FileMode
}

type clients struct {
	Data  [][]interface{}
	Hosts []string
}

// Reads info from config file
func readConfigFile(configFilePath string) {

	_, err := os.Stat(configFilePath)
	if err != nil {
		log.Fatal("Config file is missing: ", configFilePath)
	}

	var fileConfig tomlConfig
	if _, err := toml.DecodeFile(configFilePath, &fileConfig); err != nil {
		log.Fatal(err)
	}
	//fmt.Println(fileConfig)
	Settings.Storages = fileConfig.Storages
	Settings.Address.port = fileConfig.Server.Port
}
