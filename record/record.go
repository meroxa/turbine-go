package record

const (
	JSONWithSchemaFormat Format = "jsonWithSchema"
	JSONFormat                  = "json"
	RawFormat                   = "raw"
)

type Format string

type Editable interface {
	Get(path string) (Value, error)
	Set(path string, value interface{}) error
	Delete(path string) (bool, error)
	Format() Format
}

type Payload []byte
