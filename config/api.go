package config

// API is a type for API server settings.
type API struct {
	Host string `toml:"host"`
	Path string `toml:"path"`
	Port int    `toml:"port"`
}

// DefaultAPI is the default API server configuration.
var DefaultAPI = API{
	Host: "api",
	Path: "/",
	Port: 80,
}
