// Package confusion provides a simple mechanism for loading configuration
// files of arbitrary formats from a number of alternate locations.
package confusion

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/njkleiner/go-expandstrict"
	"github.com/spf13/afero"
)

var (
	// ErrConfigNotFound is returned when LoadConfig is unable to locate
	// a suitable configuration file in any of the provided locations.
	ErrConfigNotFound = errors.New("config not found")
	// ErrLoaderNotFound is returned by LoadConfig if a configuration
	// file is found but it cannot be loaded because no Loader matching
	// its file extension is available.
	ErrLoaderNotFound = errors.New("loader not found")
)

// Options determines where LoadConfig will look for configuration files
// and what formats can be loaded.
//
// Loaders maps file extensions (e.g., ".json") to a Loader
// that is capable of loading configuration files with this extension.
type Options struct {
	Prefix  string
	Paths   []string
	Loaders map[string]Loader

	fs afero.Fs
}

// A Loader is responsible for loading configuration files of a specific format.
type Loader interface {
	// Load attemtps to load a configuration file from r and store the
	// result in v.
	Load(r io.Reader, v interface{}) error
}

// LoadConfig attempts to load a configuration file with a certain name
// according to the given Options and store the result in config.
func LoadConfig(name string, opts Options, config interface{}) error {
	if opts.fs == nil {
		opts.fs = afero.NewOsFs()
	}

	for _, path := range opts.Paths {
		exp, err := expandstrict.Expand(path, os.Getenv)

		if err != nil {
			continue
		}

		path = filepath.Join(exp, opts.Prefix, name)

		info, err := opts.fs.Stat(path)

		if err != nil || info.IsDir() {
			continue
		}

		ext := filepath.Ext(info.Name())
		loader, ok := opts.Loaders[ext]

		if !ok {
			return ErrLoaderNotFound
		}

		file, err := opts.fs.Open(path)

		if err != nil {
			continue
		}

		defer file.Close()

		err = loader.Load(file, config)

		if err == nil {
			return nil
		}
	}

	return ErrConfigNotFound
}
