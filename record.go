package valve

import (
	"encoding/json"
	"time"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Records struct {
	Stream  string
	records []Record
}

func NewRecords(rr []Record) Records {
	return Records{records: rr}
}

func GetRecords(r Records) []Record {
	return r.records
}

type RecordsWithErrors struct {
	Stream  string
	records []RecordWithError
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
		return false
	}

	return false
}

type Payload []byte

func (p Payload) Map() (map[string]interface{}, error) {
	var m map[string]interface{}
	err := json.Unmarshal(p, &m)
	return m, err
}

func (p Payload) Get(path string) interface{} {
	return gjson.Get(string(p), path).Value()
}

func (p *Payload) Set(path string, value interface{}) error {
	val, err := sjson.Set(string(*p), path, value)
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
