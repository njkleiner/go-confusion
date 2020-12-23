package toml

import (
	"io"

	"github.com/pelletier/go-toml"
)

var Loader = TOMLLoader{}

type TOMLLoader struct{}

func (TOMLLoader) Load(r io.Reader, v interface{}) error {
	return toml.NewDecoder(r).Decode(v)
}
