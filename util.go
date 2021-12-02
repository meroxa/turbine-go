package valve

import (
	"math/rand"
	"time"
)

func randID() int {
	return rand.Intn(10000)
}

func wrapRecord(key string, m map[string]interface{}) Record {
	return Record{
		Key: key,
		Payload: m,
		Timestamp: time.Now(),
	}
}