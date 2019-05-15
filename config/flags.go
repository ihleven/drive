package config

import (
	"fmt"
	"os"

	arg "github.com/alexflint/go-arg"
)

// Args defines the application configuration
type args struct {
	Port           int    `arg:"env"`
	Verbose        bool   `arg:"-v" help:"verbosity level"`
	ConfigFilePath string `arg:"env:CLOUD11_CONFIG"`
}

// Version return the application name and version
func (a args) Version() string {
	return "cloud11-v0.0.1"
}

// Description explains the function of the application
func (a args) Description() string {
	return "\nthis program does this and that\n"
}

// Args holds the current application configuration defined in Args
var Args args

// ParseArgs parses the application configuration into Config by reading command line args and environment variables
func ParseArgs() settings {

	Args.ConfigFilePath = "cloud11.toml"
	arg.MustParse(&Args)

	readConfigFile(Args.ConfigFilePath)

	if Args.Port != 0 {
		Settings.Address.port = Args.Port
	}
	return Settings
}

func getCwd() string {

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error while getting current directory.")
		return ""
	}
	return cwd
}
