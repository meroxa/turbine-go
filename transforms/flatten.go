package transforms

import (
	"encoding/json"
	"github.com/meroxa/turbine-go"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"strconv"
	"strings"
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

// FlattenSub takes a potentially nested JSON payload and a path (in dot notation e.g. "foo.bar") and
// returns a JSON object with only the nested structure at the path specified flattened.
func FlattenSub(p *turbine.Payload, path string) error {
	return FlattenSubWithDelimiter(p, path, ".")
}

// FlattenSubWithDelimiter is a variant of FlattenSub that supports a custom delimiter.
func FlattenSubWithDelimiter(p *turbine.Payload, path string, del string) error {
	sub := gjson.GetBytes(*p, path)

	var child map[string]interface{}
	err := json.Unmarshal([]byte(sub.String()), &child)

	// wrap sub with parent
	hops := strings.Split(path, ".")
	lastNode := hops[len(hops)-1]
	previousNodes := strings.Join(hops[:len(hops)-1], del)
	parent := make(map[string]interface{})
	parent[lastNode] = child

	// flatten the subtree
	if err != nil {
		return err
	}
	out := make(map[string]interface{})
	flatten("\\"+del, parent, out)

	// remove the subtree from the original object
	res, err := sjson.DeleteBytes(*p, path)
	if err != nil {
		return err
	}

	// set all of the flattened keys at the correct level
	for k, v := range out {
		newPath := strings.Join([]string{previousNodes, k}, ".")
		res, err = sjson.SetBytes(res, newPath, v)
	}
	if err != nil {
		return err
	}

	*p = res
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
