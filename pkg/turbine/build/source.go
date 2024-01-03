package build

import (
	"context"

	sdk "github.com/meroxa/turbine-go/v3/pkg/turbine"

	"github.com/conduitio/conduit-commons/opencdc"
	client "github.com/meroxa/turbine-core/v2/pkg/client"
	"github.com/meroxa/turbine-core/v2/proto/turbine/v2"
)

type source struct {
	streamName string
	c          client.Client
}

func (s *source) Read() (sdk.Records, error) {
	return s.ReadWithContext(context.Background())
}

func (s *source) ReadWithContext(ctx context.Context) (sdk.Records, error) {
	resp, err := s.c.ReadRecords(ctx, &turbinev2.ReadRecordsRequest{
		SourceStream: s.streamName,
	})
	if err != nil {
		return sdk.Records{}, err
	}

	out := make([]opencdc.Record, len(resp.StreamRecords.Records))
	for i, r := range resp.StreamRecords.Records {
		if err := out[i].FromProto(r); err != nil {
			return sdk.Records{}, err
		}
	}

	return sdk.Records{
		Stream:  resp.StreamRecords.StreamName,
		Records: out,
	}, nil
}
