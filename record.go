package valve

import "time"

type Record struct {
	Key string
	Payload map[string]interface{}
	Timestamp time.Time
}

type RecordWithError struct {
	Error error
	Record
}