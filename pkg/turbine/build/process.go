package build

import (
	"context"
	"reflect"
	"strings"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	sdk "github.com/meroxa/turbine-go/v2/pkg/turbine"
)

func (b *builder) Process(rs sdk.Records, fn sdk.Function) (sdk.Records, error) {
	return b.ProcessWithContext(context.Background(), rs, fn)
}

func (b *builder) ProcessWithContext(ctx context.Context, rs sdk.Records, fn sdk.Function) (sdk.Records, error) {
	c, err := b.Client.Process(
		ctx,
		&pb.ProcessRecordsRequest{
			Process: &pb.ProcessRecordsRequest_Process{
				Name: strings.ToLower(reflect.TypeOf(fn).Name()),
			},
			Records: recordsToCollection(rs),
		})
	if err != nil {
		return sdk.Records{}, err
	}

	out := collectionToRecords(c)
	if b.runProcess {
		out.Records = fn.Process(out.Records)
	}
	return out, nil
}
