package record

import (
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type JSON struct {
	payload Payload
	format  Format
}

func (r *JSON) Get(path string) (Value, error) {
	return Value{gjson.Get(string(r.payload), path).Value()}, nil
}

func (r *JSON) Set(path string, value interface{}) error {
	_, err := sjson.Set(string(r.payload), path, value)
	return err
}

func (r *JSON) Delete(path string) (bool, error) {
	if !gjson.Get(string(r.payload), path).Exists() {
		return false, nil
	}

	// update payload
	_, err := sjson.Delete(string(r.payload), path)
	return true, err
}
