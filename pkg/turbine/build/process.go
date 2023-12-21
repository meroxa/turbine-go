package build

import (
	"context"
	"reflect"
	"strings"

	"github.com/conduitio/conduit-commons/opencdc"
	"github.com/conduitio/conduit-commons/proto/opencdc/v1"
	"github.com/meroxa/turbine-core/v2/proto/turbine/v2"
	sdk "github.com/meroxa/turbine-go/v3/pkg/turbine"
)

func (b *builder) Process(rs sdk.Records, fn sdk.Function) (sdk.Records, error) {
	return b.ProcessWithContext(context.Background(), rs, fn)
}

func (b *builder) ProcessWithContext(ctx context.Context, rs sdk.Records, fn sdk.Function) (sdk.Records, error) {
	protoRecords := make([]*opencdcv1.Record, len(rs.Records))
	for i, r := range rs.Records {
		protoRecords[i] = &opencdcv1.Record{}
		if err := r.ToProto(protoRecords[i]); err != nil {
			return sdk.Records{}, err
		}
	}

	resp, err := b.c.ProcessRecords(ctx, &turbinev2.ProcessRecordsRequest{
		StreamRecords: &turbinev2.StreamRecords{
			StreamName: rs.Stream,
			Records:    protoRecords,
		},
		Process: &turbinev2.ProcessRecordsRequest_Process{
			Name: strings.ToLower(reflect.TypeOf(fn).Name()),
		},
	})
	if err != nil {
		return sdk.Records{}, err // todo: wrap err
	}

	processedRecords := make([]opencdc.Record, len(resp.StreamRecords.Records))
	for i, r := range resp.StreamRecords.Records {
		if err := processedRecords[i].FromProto(r); err != nil {
			return sdk.Records{}, err
		}
	}

	if b.runProcess {
		return sdk.Records{
			Stream:  rs.Stream,
			Records: fn.Process(processedRecords),
		}, nil
	}

	return sdk.Records{
		Stream:  rs.Stream,
		Records: processedRecords,
	}, nil
}
