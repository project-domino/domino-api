package main

import (
	"log"
	"os"

	"github.com/project-domino/domino-go/config"
)

// Config is the configuration for the server.
var Config ConfigType

// ConfigType is the type of the configuration for the server.
type ConfigType struct {
	API       config.API       `toml:"api"`
	HTTP      config.HTTP      `toml:"http"`
	Templates config.Templates `toml:"templates"`
}

func init() {
	// Create default config object.
	Config = ConfigType{
		API:       config.DefaultAPI,
		HTTP:      config.DefaultHTTP,
		Templates: config.DefaultTemplates,
	}

	// Read config or die.
	if err := config.LoadConfig(&Config, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
