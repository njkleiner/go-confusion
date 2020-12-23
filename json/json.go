package json

import (
	"encoding/json"
	"io"
)

var Loader = JSONLoader{}

type JSONLoader struct{}

func (JSONLoader) Load(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}
