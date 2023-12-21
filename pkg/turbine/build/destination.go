package build

import (
	"context"

	sdk "github.com/meroxa/turbine-go/v3/pkg/turbine"

	"github.com/conduitio/conduit-commons/proto/opencdc/v1"
	"github.com/meroxa/turbine-core/v2/pkg/client"
	"github.com/meroxa/turbine-core/v2/proto/turbine/v2"
)

type destination struct {
	id string
	c  client.Client
}

func (d *destination) Write(rr sdk.Records) error {
	return d.WriteWithContext(context.Background(), rr)
}

func (d *destination) WriteWithContext(ctx context.Context, rr sdk.Records) error {
	protoRecords := make([]*opencdcv1.Record, len(rr.Records))
	for i, r := range rr.Records {
		protoRecords[i] = &opencdcv1.Record{}
		if err := r.ToProto(protoRecords[i]); err != nil {
			return err
		}
	}

	if _, err := d.c.WriteRecords(ctx, &turbinev2.WriteRecordsRequest{
		DestinationID: d.id,
		StreamRecords: &turbinev2.StreamRecords{
			StreamName: rr.Stream,
			Records:    protoRecords,
		},
	}); err != nil {
		return err
	}

	return nil
}
