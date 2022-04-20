package turbine

import (
	"github.com/tidwall/gjson"
	"testing"
)

func TestPayload_Set(t *testing.T) {
	type args struct {
		path  string
		value interface{}
	}
	tests := []struct {
		name    string
		p       Payload
		args    args
		wantErr bool
	}{
		{"add new", []byte(`{"schema":{"type":"struct","fields":[{"type":"int32","optional":false,"field":"id"}],"optional":false,"name":"users"},"payload":{"id":15}}`), args{"username", "test"}, false},
		{"update", []byte(`{"schema":{"type":"struct","fields":[{"type":"int32","optional":false,"field":"id"}],"optional":false,"name":"users"},"payload":{"id":15}}`), args{"id", 16}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.Set(tt.args.path, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.name == "add new" {
				schemaFields := gjson.Get(string(tt.p), "schema.fields")
				if l := len(schemaFields.Array()); l != 2 {
					t.Errorf("Set() fields len = %v, want %v", l, 2)
				}
			}
		})
	}
}
