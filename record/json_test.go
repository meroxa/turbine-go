package record

import (
	"reflect"
	"testing"
	"time"
)

func simpleJSON() string {
	return "hello"
}

func simpleJSONMap() map[string]interface{} {
	return map[string]interface{}{
		"message": "hello",
	}
}

func nestedJSONField() map[string]interface{} {
	return map[string]interface{}{
		"nested": map[string]interface{}{
			"message": "hello",
		},
	}
}

func Test_parseField(t *testing.T) {
	type args struct {
		f interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    Field
		wantErr bool
	}{
		{"simple", args{simpleJSON()}, Field{
			Type:     "string",
			Value:    "hello",
			Required: false,
			Default:  nil,
		}, false},
		{"map", args{simpleJSONMap()}, Field{
			Type:     "map",
			Value:    map[string]interface{}{"message": "hello"},
			Required: false,
			Default:  nil,
		}, false},
		{"map", args{nestedJSONField()}, Field{
			Type:     "map",
			Value:    map[string]interface{}{"nested": map[string]interface{}{"message": "hello"}},
			Required: false,
			Default:  nil,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseField(tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseField() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseField() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSON_ToInternal(t *testing.T) {
	ts := time.Now()
	type fields struct {
		key       string
		payload   Payload
		timestamp time.Time
		format    Format
	}
	tests := []struct {
		name    string
		fields  fields
		want    Internal
		wantErr bool
	}{
		{
			"simple",
			fields{"1", []byte(`{"message":"hello"}`), ts, JSONFormat},
			Internal{
				Key: "1",
				Value: map[string]Field{
					"message": Field{"string", "hello", false, nil},
				},
				Timestamp: ts,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &JSON{
				key:       tt.fields.key,
				payload:   tt.fields.payload,
				timestamp: tt.fields.timestamp,
				format:    tt.fields.format,
			}
			got, err := r.ToInternal()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToInternal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToInternal() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInternal_ToJSON(t *testing.T) {
	ts := time.Now()
	type fields struct {
		Key       string
		Value     map[string]Field
		Timestamp time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		want    JSON
		wantErr bool
	}{
		{"simple", fields{
			Key: "1",
			Value: map[string]Field{
				"message": Field{"string", "hello", false, nil},
			},
			Timestamp: ts,
		}, JSON{
			key:       "1",
			payload:   []byte(`{"message":"hello"}`),
			timestamp: ts,
			format:    JSONFormat,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Internal{
				Key:       tt.fields.Key,
				Value:     tt.fields.Value,
				Timestamp: tt.fields.Timestamp,
			}
			got, err := i.ToJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}
