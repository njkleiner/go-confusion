package confusion

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

var (
	ErrConfigNotFound = errors.New("config not found")
	ErrLoaderNotFound = errors.New("loader not found")
)

type Options struct {
	Prefix      string
	UserPaths   []string
	SystemPaths []string
	Loaders     map[string]Loader

	fs afero.Fs
}

type Loader interface {
	Load(r io.Reader, v interface{}) error
}

func LoadConfig(name string, opts Options, config interface{}) error {
	if opts.fs == nil {
		opts.fs = afero.NewOsFs()
	}

	for _, path := range append(opts.UserPaths, opts.SystemPaths...) {
		path = filepath.Join(os.ExpandEnv(path), opts.Prefix, name)

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
