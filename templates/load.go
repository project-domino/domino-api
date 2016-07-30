package templates

import (
	"archive/zip"
	"html/template"
	"io/ioutil"
	"path"

	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/zipfs"

	"github.com/project-domino/domino-go/config"
)

// Load loads templates based on the given configuration.
func Load(templateConfig config.Templates) (*template.Template, error) {
	fs, err := getFileSystem(templateConfig)
	if err != nil {
		return nil, err
	}

	allFiles, err := fs.ReadDir("/")
	if err != nil {
		return nil, err
	}

	t := template.New("").Funcs(template.FuncMap{
		"toSnakeCase": ToSnakeCase,
	})
	for _, file := range allFiles {
		if file.IsDir() {
			continue
		}

		reader, err := fs.Open("/" + file.Name())
		if err != nil {
			return nil, err
		}

		src, err := ioutil.ReadAll(reader)
		if err != nil {
			return nil, err
		}

		t.New(file.Name()).Parse(string(src))
	}
	return t, nil
}

func getFileSystem(templateConfig config.Templates) (vfs.FileSystem, error) {
	if templateConfig.Dev {
		return vfs.OS(templateConfig.Path), nil
	}

	reader, err := zip.OpenReader(templateConfig.Path)
	if err != nil {
		return nil, err
	}
	return zipfs.New(reader, path.Base(templateConfig.Path)), nil
}
