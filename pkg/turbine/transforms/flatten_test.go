package transforms

import (
	"log"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/meroxa/turbine-go/v2/pkg/turbine"
)

func TestFlatten(t *testing.T) {
	tests := []struct {
		name    string
		p       turbine.Payload
		want    turbine.Payload
		wantErr bool
	}{
		{"simple", []byte(`{"user": {"id":16, "name": "alice"}}`), []byte(`{"user.id":16,"user.name":"alice"}`), false},
		{"arrays", []byte(`{"user": {"locations":[{"city":"London", "country":"UK"},{"city":"San Francisco","country":"USA"}]}}`), []byte(`{"user.locations.0.city":"London","user.locations.0.country":"UK","user.locations.1.city":"San Francisco","user.locations.1.country":"USA"}`), false},
		{"non-cdc record", []byte(`{"id": 1, "user": {"id": 100, "name": "alice", "email": "alice@example.com"}, "actions": ["register", "purchase"]}`), []byte(`{"actions.0":"register","actions.1":"purchase","id":1,"user.email":"alice@example.com","user.id":100,"user.name":"alice"}`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Flatten(&tt.p); (err != nil) != tt.wantErr {
				t.Errorf("Flatten() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.p, tt.want) {
				log.Printf("p: %+v", string(tt.p))
				t.Errorf("Flatten() got = %v, want %v", string(tt.p), string(tt.want))
			}
		})
	}
}

func TestFlattenSub(t *testing.T) {
	type args struct {
		p    turbine.Payload
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    turbine.Payload
		wantErr bool
	}{
		{"nested", args{[]byte(`{"user":{"id":16,"name":"alice","nested":{"id":1,"message":"hello"}}}`), "user.nested"}, []byte(`{"user":{"id":16,"name":"alice","nested.message":"hello","nested.id":1}}`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := FlattenSub(&tt.args.p, tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("FlattenSub() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.JSONEq(t, string(tt.args.p[:]), string(tt.want[:]))
		})
	}
}

func TestFlattenSubWithDelimiter(t *testing.T) {
	type args struct {
		p    turbine.Payload
		path string
		del  string
	}
	tests := []struct {
		name    string
		args    args
		want    turbine.Payload
		wantErr bool
	}{
		{"custom delimiter", args{[]byte(`{"user":{"id":16,"name":"alice","nested":{"id":1,"message":"hello"}}}`), "user.nested", "-"}, []byte(`{"user":{"id":16,"name":"alice","nested-id":1,"nested-message":"hello"}}`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := FlattenSubWithDelimiter(&tt.args.p, tt.args.path, tt.args.del); (err != nil) != tt.wantErr {
				t.Errorf("FlattenSub() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.JSONEq(t, string(tt.args.p[:]), string(tt.want[:]))
		})
	}
}
