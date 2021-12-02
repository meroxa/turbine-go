package valve

import "time"

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