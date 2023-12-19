package build

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/meroxa/turbine-core/v2/lib/go/github.com/meroxa/turbine/core"
	sdk "github.com/meroxa/turbine-go/v3/pkg/turbine"
)

func fromRecords(rs sdk.Records) *pb.StreamRecords {
	rr := make([]*pb.Record, len(rs.Records))

	for i, r := range rs.Records {
		rr[i] = &pb.Record{
			Key:       r.Key,
			Value:     r.Payload,
			Timestamp: timestamppb.New(r.Timestamp),
		}
	}

	return &pb.StreamRecords{
		StreamName: rs.Stream,
		Records:    rr,
	}
}

func toRecords(sr *pb.StreamRecords) sdk.Records {
	rs := make([]sdk.Record, len(sr.Records))

	for i, r := range sr.Records {
		rs[i] = sdk.Record{
			Key:       r.Key,
			Payload:   r.Value,
			Timestamp: r.Timestamp.AsTime(),
		}
	}

	return sdk.Records{
		Stream:  sr.StreamName,
		Records: rs,
	}
}
