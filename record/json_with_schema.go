package record

import (
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"strconv"
	"strings"
)

type JSONWithSchema struct {
	payload Payload
	format  Format
}

type schemaField struct {
	Field    string `json:"field"`
	Optional bool   `json:"optional"`
	Type     string `json:"type"`
}

func (r *JSONWithSchema) Get(path string) (Value, error) {
	nestedPath := strings.Join([]string{"payload", path}, ".")
	return Value{gjson.Get(string(r.payload), nestedPath).Value()}, nil
}

func (r *JSONWithSchema) Set(path string, value interface{}) error {
	nestedPath := strings.Join([]string{"payload", path}, ".")
	fieldExists := gjson.Get(string(r.payload), nestedPath).Exists()

	// update payload
	val, err := sjson.Set(string(r.payload), nestedPath, value)
	if err != nil {
		return err
	}
	r.payload = []byte(val)

	// Add schema field if field is new
	if !fieldExists {
		fieldType := mapGoToKCDataTypes(val)

		field := schemaField{
			Field:    path,
			Optional: true,
			Type:     fieldType,
		}

		schemaNestedPath := strings.Join([]string{"schema", "fields.-1"}, ".")
		sval, err := sjson.Set(string(r.payload), schemaNestedPath, field)
		if err != nil {
			return err
		}
		r.payload = []byte(sval)
	}

	return nil
}

func (r *JSONWithSchema) Delete(path string) (bool, error) {
	nestedPath := strings.Join([]string{"payload", path}, ".")
	if !gjson.Get(string(r.payload), nestedPath).Exists() {
		return false, nil
	}

	// update payload
	val, err := sjson.Delete(string(r.payload), nestedPath)
	if err != nil {
		return true, err
	}

	// Remove schema field
	schemaFields := gjson.Get(val, "schema.fields").Array()
	for i, field := range schemaFields {
		if field.Map()["field"].String() == path {
			fieldPath := strings.Join([]string{"schema.fields", strconv.Itoa(i)}, ".")
			val, err = sjson.Delete(val, fieldPath)
			if err != nil {
				return true, err
			}
			break
		}
	}

	r.payload = []byte(val)
	return true, nil
}

func (r *JSONWithSchema) Format() Format {
	return r.format
}

// map Go types to Apache Kafka Connect data types
func mapGoToKCDataTypes(v interface{}) string {
	switch v.(type) {
	case string:
		return "string"
	case int8:
		return "int8"
	case int16:
		return "int16"
	case int, int32:
		return "int32"
	case int64:
		return "int64"
	case float32:
		return "float32"
	case float64:
		return "float64"
	case bool:
		return "boolean"
	default:
		return "unsupported"
	}
}
