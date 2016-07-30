package config

// Templates is a type for template configuration settings.
type Templates struct {
	Dev  bool   `toml:"dev"`
	Path string `toml:"path"`
}

// DefaultTemplates is the default template configuration.
var DefaultTemplates = Templates{
	Dev:  false,
	Path: "templates.zip",
}
