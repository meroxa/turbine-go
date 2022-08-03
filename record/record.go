package record

import "time"

const (
	JSONWithSchemaFormat Format = "jsonWithSchema"
	JSONFormat                  = "json"
	RawFormat                   = "raw"
)

type Format string

type Record interface {
	Get(path string) (Value, error)
	Set(path string, value interface{}) error
	Delete(path string) (bool, error)
	Key() string
	Timestamp() time.Time
	Format() Format
}

type Payload []byte

func DetectFormat(raw []byte) Format {
	return RawFormat
}
