package build

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	sdk "github.com/meroxa/turbine-go/v2/pkg/turbine"
)

func recordsToCollection(rs sdk.Records) *pb.Collection {
	rds := []*pb.Record{}
	for _, r := range rs.Records {
		rds = append(rds,
			&pb.Record{
				Key:       r.Key,
				Value:     r.Payload,
				Timestamp: timestamppb.New(r.Timestamp),
			})
	}
	return &pb.Collection{
		Stream:  rs.Stream,
		Records: rds,
		Name:    rs.Name,
	}
}

func collectionToRecords(c *pb.Collection) sdk.Records {
	rs := []sdk.Record{}
	for _, r := range c.Records {
		rs = append(rs,
			sdk.Record{
				Key:       r.Key,
				Payload:   r.Value,
				Timestamp: r.Timestamp.AsTime(),
			},
		)
	}

	return sdk.Records{
		Stream:  c.Stream,
		Records: rs,
		Name:    c.Name,
	}
}
