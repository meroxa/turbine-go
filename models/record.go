package models

import (
	"encoding/json"
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"log"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Records struct {
	Stream  string
	records []Record
	pcol    beam.PCollection
}

func NewRecords(rr []Record) Records {
	return Records{records: rr}
}

func NewRecordsWithPcol(pcol beam.PCollection) Records {
	return Records{pcol: pcol}
}

func GetRecords(r Records) []Record {
	return r.records
}

func GetPcol(r Records) beam.PCollection {
	return r.pcol
}

type Record struct {
	Key       string
	Payload   Payload
	Timestamp time.Time
}

// JSONSchema returns true if the record is formatted with JSON Schema, false otherwise
func (r Record) JSONSchema() bool {
	p, err := r.Payload.Map()
	if err != nil {
		return false
	}

	if _, ok := p["schema"]; ok {
		if _, ok := p["payload"]; ok {
			return true
		}
	}

	return false
}

// OpenCDC returns true if the record is formatted with OpenCDC schema, false otherwise
func (r Record) OpenCDC() bool {
	p, err := r.Payload.Map()
	if err != nil {
		return false
	}

	if _, ok := p["schema"]; ok {
		if payload, ok := p["payload"]; ok {
			if _, ok := payload.(map[string]interface{})["after"]; ok {
				return true
			}
		}
	}

	return false
}

func (r Record) String() string {
	b, err := json.Marshal(r.Payload)
	if err != nil {
		log.Printf("error marshaling to string; err: %s", err.Error())
	}

	return string(b)
}

func (r Record) MarshalJSON() ([]byte, error) {
	log.Println("MarshalJSON called")
	return nil, nil
}

func (r Record) UnmarshalJSON(buf []byte) error {
	log.Println("UnmarshalJSON called")
	return nil
}

type Payload []byte

func (p Payload) Map() (map[string]interface{}, error) {
	var m map[string]interface{}
	err := json.Unmarshal(p, &m)
	return m, err
}

func (p Payload) Get(path string) interface{} {
	nestedPath := strings.Join([]string{"payload", path}, ".")
	return gjson.Get(string(p), nestedPath).Value()
}

// TODO: Add GetType(path string) to tell you what the data type is.
// TODO: Should we passthrough the gjson helper methods?

type schemaField struct {
	Field    string `json:"field"`
	Optional bool   `json:"optional"`
	Type     string `json:"type"`
}

func (p *Payload) Set(path string, value interface{}) error {
	nestedPath := strings.Join([]string{"payload", path}, ".")
	fieldExists := gjson.Get(string(*p), nestedPath).Exists()

	// update payload
	val, err := sjson.Set(string(*p), nestedPath, value)
	if err != nil {
		return err
	}
	*p = []byte(val)

	// Add schema field if field is new
	if !fieldExists {
		fieldType := mapGoToKCDataTypes(val)

		field := schemaField{
			Field:    path,
			Optional: true,
			Type:     fieldType,
		}

		schemaNestedPath := strings.Join([]string{"schema", "fields.-1"}, ".")
		sval, err := sjson.Set(string(*p), schemaNestedPath, field)
		if err != nil {
			return err
		}
		*p = []byte(sval)
	}

	return nil
}

func (p *Payload) Delete(path string) error {
	val, err := sjson.Delete(string(*p), path)
	if err != nil {
		return err
	}
	*p = []byte(val)
	return nil
}

type RecordWithError struct {
	Error error
	Record
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
