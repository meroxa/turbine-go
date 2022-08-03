package record

import (
	"encoding/json"
	"time"
)

type Internal struct {
	Key       string
	Value     map[string]Field
	Timestamp time.Time
}

type FieldType string

type Field struct {
	Type     FieldType
	Value    interface{}
	Required bool
	Default  interface{}
}

func (f Field) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.Value)
}

func (f Field) Schema(name string) schemaField {
	return schemaField{
		Field:    name,
		Optional: !f.Required,
		Type:     mapGoToKCDataTypes(f.Value),
	}
}

func (i *Internal) ToJSONWithSchema() (JSONWithSchema, error) {
	return JSONWithSchema{}, nil
}
