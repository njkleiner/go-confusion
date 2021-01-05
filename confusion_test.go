package confusion

import (
	"io"
	"os"
	"testing"

	"github.com/njkleiner/go-confusion/json"
	"github.com/njkleiner/go-confusion/toml"
	"github.com/spf13/afero"
)

// TestLoadConfig is a basic sanity test.
func TestLoadConfig(t *testing.T) {
	opts := makeOptions()

	opts.fs.MkdirAll("/etc/confusion", os.ModePerm)
	write(opts.fs, "/etc/confusion/config.toml", "message = \"Hello World\"\r\n")

	type config struct {
		Message string
	}

	c := config{}

	err := LoadConfig("config.toml", opts, &c)

	if err != nil {
		t.Error(err)
	}

	if c.Message != "Hello World" {
		t.Errorf("invalid config value: expected message = Hello World, actual message = %s", c.Message)
	}
}

// TestLoadConfigPaths tests loading the config from a number of
// alternate paths.
func TestLoadConfigPaths(t *testing.T) {
	opts := makeOptions()

	os.Setenv("XDG_CONFIG_HOME", "")
	os.Setenv("HOME", "/home/test")

	opts.fs.MkdirAll("/home/test/.config/confusion", os.ModePerm)
	write(opts.fs, "/home/test/.config/confusion/config.toml", "message = \"Hello World\"\r\n")

	type config struct {
		Message string
	}

	c := config{}

	err := LoadConfig("config.toml", opts, &c)

	if err != nil {
		t.Error(err)
	}

	if c.Message != "Hello World" {
		t.Errorf("invalid config value: expected message = Hello World, actual message = %s", c.Message)
	}
}

// TestLoadConfigErrLoaderNotFound exists mainly for coverage reasons.
func TestLoadConfigErrLoaderNotFound(t *testing.T) {
	opts := makeOptions()
	opts.Loaders = make(map[string]Loader)

	opts.fs.MkdirAll("/etc/confusion", os.ModePerm)
	write(opts.fs, "/etc/confusion/config.toml", "message = \"Hello World\"\r\n")

	type config struct {
		Message string
	}

	c := config{}

	err := LoadConfig("config.toml", opts, &c)

	if err != ErrLoaderNotFound {
		t.Errorf("unexpected error; expected ErrLoaderNotFound, actual %v", err)
	}
}

// TestLoadConfigErrConfigNotFound exists mainly for coverage reasons.
func TestLoadConfigErrConfigNotFound(t *testing.T) {
	opts := makeOptions()

	opts.fs.MkdirAll("/etc/confusion", os.ModePerm)

	type config struct {
		Message string
	}

	c := config{}

	err := LoadConfig("config.toml", opts, &c)

	if err != ErrConfigNotFound {
		t.Errorf("unexpected error; expected ErrConfigNotFound, actual %v", err)
	}
}

func write(fs afero.Fs, dest, data string) {
	file, err := fs.OpenFile(dest, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	io.WriteString(file, data)
}

func makeOptions() Options {
	loaders := make(map[string]Loader)
	loaders[".json"] = json.Loader
	loaders[".toml"] = toml.Loader

	return Options{
		Prefix:      "confusion",
		UserPaths:   []string{"$XDG_CONFIG_HOME", "$HOME/.config"},
		SystemPaths: []string{"/etc"},
		Loaders:     loaders,
		fs:          afero.NewMemMapFs(),
	}
}
