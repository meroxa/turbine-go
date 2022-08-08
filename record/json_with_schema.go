package record

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"strconv"
	"strings"
	"time"
)

type JSONWithSchema struct {
	key       string
	payload   Payload
	timestamp time.Time
	format    Format
}

type schemaField struct {
	Field    string `json:"field"`
	Optional bool   `json:"optional"`
	Type     string `json:"type"`
	Default  string `json:"default,omitempty"`
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

func (r *JSONWithSchema) Key() string {
	return r.key
}

func (r *JSONWithSchema) Timestamp() time.Time {
	return r.timestamp
}

func (r *JSONWithSchema) Format() Format {
	return r.format
}

func (r *JSONWithSchema) ToInternal() (Internal, error) {
	var unstructured struct {
		Key   string `json:"key"`
		Value struct {
			Schema  map[string]schemaField `json:"schema"`
			Payload map[string]interface{} `json:"payload"`
		} `json:"value"`
		Timestamp string `json:"timestamp"`
	}
	err := json.Unmarshal(r.payload, &unstructured)
	if err != nil {
		return Internal{}, err
	}

	value := make(map[string]Field)
	for n, f := range unstructured.Value.Payload {
		field := Field{
			Type:     FieldType(mapKCDataTypesToGo(unstructured.Value.Schema[n].Type)),
			Value:    f,
			Required: !unstructured.Value.Schema[n].Optional,
			Default:  unstructured.Value.Schema[n].Default,
		}
		value[n] = field
	}
	return Internal{
		Key:       r.Key(),
		Value:     value,
		Timestamp: r.Timestamp(),
	}, nil
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

// map Apache Kafka Connect data types to Go types
func mapKCDataTypesToGo(v string) string {
	switch v {
	case "boolean":
		return "bool"
	default:
		return v
	}
}
