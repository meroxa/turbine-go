package valve

import "time"

type Records struct {
	Stream string
	records []Record
}

func NewRecords(rr []Record) Records {
	return Records{records: rr}
}

type RecordsWithErrors struct {
	Stream string
	records []RecordWithError
}

type Record struct {
	Key string
	Payload Payload
	Timestamp time.Time
}

type Payload map[string]interface{}

type RecordWithError struct {
	Error error
	Record
}