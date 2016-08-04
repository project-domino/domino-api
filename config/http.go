package config

import "fmt"

// HTTP is a type for HTTP server settings.
type HTTP struct {
	Debug bool   `toml:"debug"`
	Host  string `toml:"host,omitempty"`
	Port  int    `toml:"port"`
}

// DefaultHTTP is the default HTTP server configuration.
var DefaultHTTP = HTTP{
	Debug: false,
	Host:  "",
	Port:  80,
}

// ServeOn returns a string appropriate for e.g. http.ListenAndServe.
func (config HTTP) ServeOn() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
