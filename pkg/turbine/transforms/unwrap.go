package transforms

import (
	"encoding/json"

	"github.com/meroxa/turbine-go/v2/pkg/turbine"
)

// Unwrap takes a JSON payload that may or may not be of the JSON with Schema format and returns only the "payload". In
// the case where there is no schema envelope, the original unmodified record is returned without error.
func Unwrap(p *turbine.Payload) error {
	var child map[string]interface{}
	err := json.Unmarshal(*p, &child)
	if err != nil {
		return err
	}

	if payload, ok := child["payload"]; ok {
		b, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		*p = b
	}
	return nil
}
