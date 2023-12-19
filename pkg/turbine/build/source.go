package build

import (
	"context"

	sdk "github.com/meroxa/turbine-go/v3/pkg/turbine"

	pb "github.com/meroxa/turbine-core/v2/lib/go/github.com/meroxa/turbine/core"
	client "github.com/meroxa/turbine-core/v2/pkg/client"
)

type source struct {
	streamName string
	c          client.Client
}

func (s *source) Read() (sdk.Records, error) {
	return s.ReadWithContext(context.Background())
}

func (s *source) ReadWithContext(ctx context.Context) (sdk.Records, error) {
	resp, err := s.c.ReadRecords(ctx, &pb.ReadRecordsRequest{
		SourceStream: s.streamName,
	})
	if err != nil {
		return sdk.Records{}, err
	}

	return toRecords(resp.StreamRecords), nil
}
