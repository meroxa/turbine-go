package transforms

import (
	"encoding/json"
	"strings"

	"github.com/jeremywohl/flatten"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"github.com/meroxa/turbine-go/v2/pkg/turbine"
)

// Flatten takes a potentially nested JSON payload and returns a flattened representation, using a "."
// as a delimiter. e.g. {"user": {"id":16, "name": "alice"}} becomes {"user.id":16,"user.name":"alice"}
// If an array of nested objects is encountered, the index of the element will be appended to the field
// name. e.g. {"user.location.1":"London, UK", "user.location.2":"San Francisco, USA"}
func Flatten(p *turbine.Payload) error {
	f, err := flatten.FlattenString(string(*p), "", flatten.DotStyle)
	if err != nil {
		return err
	}
	*p = []byte(f)
	return nil
}

// FlattenWithDelimiter is a variant of Flatten that supports a custom delimiter.
func FlattenWithDelimiter(p *turbine.Payload, del string) error {
	sep := flatten.SeparatorStyle{Middle: del}
	f, err := flatten.FlattenString(string(*p), "", sep)
	if err != nil {
		return err
	}
	*p = []byte(f)
	return nil
}

// FlattenSub takes a potentially nested JSON payload and a path (in dot notation e.g. "foo.bar") and
// returns a JSON object with only the nested structure at the path specified flattened.
func FlattenSub(p *turbine.Payload, path string) error {
	return FlattenSubWithDelimiter(p, path, ".")
}

// FlattenSubWithDelimiter is a variant of FlattenSub that supports a custom delimiter.
func FlattenSubWithDelimiter(p *turbine.Payload, path string, del string) error {
	sep := flatten.SeparatorStyle{Middle: del}
	sub := gjson.GetBytes(*p, path)

	var child map[string]interface{}
	err := json.Unmarshal([]byte(sub.String()), &child)
	if err != nil {
		return err
	}

	// wrap sub with parent
	hops := strings.Split(path, ".")
	lastNode := hops[len(hops)-1]
	previousNodes := strings.Join(hops[:len(hops)-1], del)
	parent := make(map[string]interface{})
	parent[lastNode] = child

	// flatten the subtree
	f, err := flatten.Flatten(parent, "", sep)
	if err != nil {
		return err
	}

	// remove the subtree from the original object
	res, err := sjson.DeleteBytes(*p, path)
	if err != nil {
		return err
	}

	// set all of the flattened keys at the correct level
	for k, v := range f {
		epath := strings.Replace(k, ".", `\.`, 1)
		newPath := strings.Join([]string{previousNodes, epath}, ".")
		res, err = sjson.SetBytes(res, newPath, v)
	}
	if err != nil {
		return err
	}

	*p = res
	return nil
}
