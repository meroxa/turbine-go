package server

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	pb "github.com/meroxa/turbine-go/v2/proto"
)

func Test_function_Process(t *testing.T) {
	var (
		proc = processor{}
		f    = function{process: proc.Process}
	)

	in := []*pb.Record{
		{
			Key:       "key-1",
			Value:     "val-1",
			Timestamp: time.Now().Unix(),
		},
		{
			Key:       "key-2",
			Value:     "val-2",
			Timestamp: time.Now().Add(1 * time.Minute).Unix(),
		},
	}

	expect := []*pb.Record{
		{
			Key:       "key-1+processed",
			Value:     "val-1",
			Timestamp: time.Now().Unix(),
		},
		{
			Key:       "key-2+processed",
			Value:     "val-2",
			Timestamp: time.Now().Add(1 * time.Minute).Unix(),
		},
	}

	resp, err := f.Process(context.Background(), &pb.ProcessRecordRequest{Records: in})
	if err != nil {
		t.Fatalf("unexpected error")
	}

	if diff := cmp.Diff(expect, resp.Records, protocmp.Transform()); diff != "" {
		t.Fatalf("mismatch (-want,+got): %s", diff)
	}
}
