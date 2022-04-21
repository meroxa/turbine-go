package transforms_test

import (
	"github.com/meroxa/turbine-go"
	"github.com/meroxa/turbine-go/transforms"
	"log"
	"reflect"
	"testing"
)

func TestPayload_Flatten(t *testing.T) {
	tests := []struct {
		name    string
		p       turbine.Payload
		want    turbine.Payload
		wantErr bool
	}{
		{"simple", []byte(`{"user": {"id":16, "name": "alice"}}`), []byte(`{"user.id":16,"user.name":"alice"}`), false},
		{"arrays", []byte(`{"user": {"locations":[{"city":"London", "country":"UK"},{"city":"San Francisco","country":"USA"}]}}`), []byte(`{"user.locations.0":{"city":"London","country":"UK"},"user.locations.1":{"city":"San Francisco","country":"USA"}}`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := transforms.Flatten(&tt.p); (err != nil) != tt.wantErr {
				t.Errorf("Flatten() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.p, tt.want) {
				log.Printf("p: %+v", string(tt.p))
				t.Errorf("Flatten() got = %v, want %v", string(tt.p), string(tt.want))
			}
		})
	}
}
