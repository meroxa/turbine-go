package valve

import "time"

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

type Payload []byte

type RecordWithError struct {
	Error error
	Record
}
