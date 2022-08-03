package record

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"log"
	"time"
)

type JSON struct {
	key       string
	payload   Payload
	timestamp time.Time
	format    Format
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

func (r *JSON) Key() string {
	return r.key
}

func (r *JSON) Timestamp() time.Time {
	return r.timestamp
}

func (r *JSON) Format() Format {
	return r.format
}

func (r *JSON) ToInternal() (Internal, error) {
	var unstructured map[string]interface{}
	err := json.Unmarshal(r.payload, &unstructured)
	if err != nil {
		return Internal{}, err
	}

	value := make(map[string]Field)
	for n, f := range unstructured {
		field, err := parseField(f)
		if err != nil {
			return Internal{}, err
		}
		value[n] = field
	}
	return Internal{
		Key:       r.Key(),
		Value:     value,
		Timestamp: r.Timestamp(),
	}, nil
}

func (i *Internal) ToJSON() (JSON, error) {
	jsonPayload, err := json.Marshal(i.Value)
	if err != nil {
		return JSON{}, err
	}
	return JSON{
		key:       i.Key,
		payload:   jsonPayload,
		timestamp: i.Timestamp,
		format:    JSONFormat,
	}, nil
}

func parseField(f interface{}) (Field, error) {
	switch f.(type) {
	case string:
		return Field{
			Type:     "string",
			Value:    f,
			Required: false,
			Default:  nil,
		}, nil
	case map[string]interface{}:
		tmp := make(map[string]Field)
		for n, f := range f.(map[string]interface{}) {
			field, err := parseField(f)
			if err != nil {
				return Field{}, err
			}
			tmp[n] = field
		}
		return Field{
			Type:     "map",
			Value:    f,
			Required: false,
			Default:  nil,
		}, nil
	default:
		log.Printf("unsupported type: %+v", f)
	}
	return Field{}, nil
}
