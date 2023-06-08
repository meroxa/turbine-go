package server

import (
	"context"
	"time"

	sdk "github.com/meroxa/turbine-go/v2/pkg/turbine"
	pb "github.com/meroxa/turbine-go/v2/proto"
)

type function struct {
	process func([]sdk.Record) []sdk.Record
}

func (f *function) Process(ctx context.Context, req *pb.ProcessRecordRequest) (*pb.ProcessRecordResponse, error) {
	var (
		in  []sdk.Record
		out []*pb.Record
	)

	for _, r := range req.Records {
		in = append(in, sdk.Record{
			Key:       r.Key,
			Payload:   sdk.Payload(r.Value),
			Timestamp: time.Unix(r.Timestamp, 0),
		})
	}

	for _, r := range f.process(in) {
		out = append(out, &pb.Record{
			Key:       r.Key,
			Value:     string(r.Payload),
			Timestamp: r.Timestamp.Unix(),
		})
	}

	return &pb.ProcessRecordResponse{Records: out}, nil
}
