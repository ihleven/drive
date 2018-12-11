package config

type DatabaseConnection struct {
	Name string
	Host string
	Port int
	User string
	Password string
}

var Dbconn = &DatabaseConnection{Name:"webcc-local", Host:"127.0.0.1", Port:1433, User:"webrx", Password:"0Q2u09KnbnawxEDEtwox"}