package build

import (
	"context"
	"reflect"
	"strings"

	pb "github.com/meroxa/turbine-core/v2/lib/go/github.com/meroxa/turbine/core"
	sdk "github.com/meroxa/turbine-go/v3/pkg/turbine"
)

func (b *builder) Process(rs sdk.Records, fn sdk.Function) (sdk.Records, error) {
	return b.ProcessWithContext(context.Background(), rs, fn)
}

func (b *builder) ProcessWithContext(ctx context.Context, rs sdk.Records, fn sdk.Function) (sdk.Records, error) {
	resp, err := b.c.ProcessRecords(ctx, &pb.ProcessRecordsRequest{
		StreamRecords: fromRecords(rs),
		Process: &pb.ProcessRecordsRequest_Process{
			Name: strings.ToLower(reflect.TypeOf(fn).Name()),
		},
	})
	if err != nil {
		return sdk.Records{}, err // todo: wrap err
	}

	rr := toRecords(resp.StreamRecords)

	if b.runProcess {
		return sdk.Records{
			Stream:  rr.Stream,
			Records: fn.Process(rr.Records),
		}, nil
	}

	return rr, nil
}
