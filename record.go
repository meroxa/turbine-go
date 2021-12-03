package valve

import "time"

type Records struct {
	records []Record
}

func NewRecords(rr []Record) Records {
	return Records{rr}
}

type RecordsWithErrors struct {
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