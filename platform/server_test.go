package platform

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/meroxa/turbine-go"
	"github.com/meroxa/turbine-go/proto"
)

func Test_protoRecordToValveRecord(t *testing.T) {
	now := time.Now()
	type args struct {
		req *proto.ProcessRecordRequest
	}
	tests := []struct {
		name string
		args args
		want []turbine.Record
	}{
		{
			name: "valid",
			args: args{
				req: &proto.ProcessRecordRequest{
					Records: []*proto.Record{
						{
							Key:       "1",
							Value:     "{ \"id\": \"2\", \"user_id\": \"100\", \"email\": \"user@example.com\", \"action\": \"logged in\" }\n",
							Timestamp: now.Unix(),
						},
					},
				},
			},
			want: []turbine.Record{
				{
					Key:       "1",
					Payload:   turbine.Payload("{ \"id\": \"2\", \"user_id\": \"100\", \"email\": \"user@example.com\", \"action\": \"logged in\" }\n"),
					Timestamp: time.Unix(now.Unix(), 0),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := protoRecordToValveRecord(tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("protoRecordToValveRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_wrapFrameworkFunc(t *testing.T) {
	pRecordsReq := &proto.ProcessRecordRequest{
		Records: []*proto.Record{
			{
				Key:       "1",
				Value:     "{ \"id\": \"2\", \"user_id\": \"100\", \"email\": \"user@example.com\", \"action\": \"logged in\" }",
				Timestamp: time.Now().Unix(),
			},
		},
	}
	vRecords := []turbine.Record{
		{
			Key:       "1",
			Payload:   turbine.Payload("{ \"id\": \"2\", \"user_id\": \"100\", \"email\": \"user@example.com\", \"action\": \"logged in\" }"),
			Timestamp: time.Now(),
		},
	}

	passthrough := func(rr []turbine.Record) []turbine.Record {
		return rr
	}

	wfn := wrapFrameworkFunc(passthrough)
	fn := struct{ ProtoWrapper }{}
	fn.ProcessMethod = wfn

	resp, err := fn.Process(context.Background(), pRecordsReq)
	if err != nil {
		t.Errorf("no error expected; got %s", err.Error())
	}

	if resp.Records == nil {
		t.Error("expect records; got nil")
	}

	if resp.Records[0].Key != vRecords[0].Key {
		t.Errorf("want key %s; got key %s", vRecords[0].Key, resp.Records[0].Key)
	}

	if resp.Records[0].Value != string(vRecords[0].Payload) {
		t.Errorf("want key %s; got key %s", string(vRecords[0].Payload), resp.Records[0].Value)
	}
}
