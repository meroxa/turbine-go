package transforms

import (
	"encoding/json"
	"github.com/meroxa/turbine-go"
	"strconv"
)

// Flatten takes a potentially nested JSON payload and returns a flattened representation, using a "."
// as a delimiter. e.g. {"user": {"id":16, "name": "alice"}} becomes {"user.id":16,"user.name":"alice"}
// If an array of nested objects is encountered, the index of the element will be appended to the field
// name. e.g. {"user.location.1":"London, UK", "user.location.2":"San Francisco, USA"}
func Flatten(p *turbine.Payload) error {
	return FlattenWithDelimiter(p, ".")
}

// FlattenWithDelimiter is a variant of Flatten that supports a custom delimiter.
func FlattenWithDelimiter(p *turbine.Payload, del string) error {
	var child map[string]interface{}
	err := json.Unmarshal(*p, &child)
	if err != nil {
		return err
	}
	out := make(map[string]interface{})
	flatten(del, child, out)

	b, err := json.Marshal(out)
	if err != nil {
		return err
	}

	*p = b
	return nil
}

func flatten(del string, src map[string]interface{}, dest map[string]interface{}) {
	for k, v := range src {
		switch child := v.(type) {
		case map[string]interface{}:
			flatten(k+del, child, dest)
		case []interface{}:
			for i := 0; i < len(child); i++ {
				dest[del+k+"."+strconv.Itoa(i)] = child[i]
			}
		default:
			dest[del+k] = v
		}
	}
}
