package turbine_test

import (
	"log"
	"reflect"
	"testing"

	"github.com/meroxa/turbine-go"
	"github.com/tidwall/gjson"
)

func TestPayload_Set(t *testing.T) {
	type args struct {
		path  string
		value interface{}
	}
	tests := []struct {
		name            string
		p               turbine.Payload
		args            args
		wantErr         bool
		schemaFieldsNum int
	}{
		{"add new", recWithSchema(), args{"email", "alice@example.com"}, false, 3},
		{"update", recWithSchema(), args{"id", 16}, false, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.Set(tt.args.path, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}

			schemaFields := gjson.Get(string(tt.p), "schema.fields")
			if l := len(schemaFields.Array()); l != tt.schemaFieldsNum {
				log.Printf("p: %+v", string(tt.p))
				t.Errorf("Set() fields len = %v, want %v", l, tt.schemaFieldsNum)
			}
		})
	}
}

func recWithSchema() turbine.Payload {
	return []byte(`
{
	"schema": {
		"type": "struct",
		"fields": [{
			"type": "int32",
			"optional": false,
			"field": "id"
		}, {
			"type": "string",
			"optional": false,
			"field": "username"
		}],
		"optional": false,
		"name": "users"
	},
	"payload": {
		"id": 15,
		"username": "test"
	}
}
`)
}

func TestPayload_Delete(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		p       turbine.Payload
		args    args
		want    turbine.Payload
		wantErr bool
	}{
		{"delete existing", []byte(`{"user":{"id":16,"name": "alice"}}`), args{"user.name"}, []byte(`{"user":{"id":16}}`), false},
		{"delete non-existent", []byte(`{"user":{"id":16,"name": "alice"}}`), args{"user.email"}, []byte(`{"user":{"id":16,"name": "alice"}}`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.Delete(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.p, tt.want) {
				// log.Printf("p: %+v", string(tt.p))
				t.Errorf("Delete() got = %v, want %v", string(tt.p), string(tt.want))
			}
		})
	}
}
