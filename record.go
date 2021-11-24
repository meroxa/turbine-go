package valve

import "time"

type Record struct {
	Key string
	Payload []byte
	Timestamp time.Time
}

type RecordWithError struct {
	Error error
	Record
}