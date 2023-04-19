package turbine

import (
	"context"
	"reflect"
	"strings"

	"github.com/meroxa/turbine-go/pkg/app"
	"github.com/meroxa/turbine-go/pkg/proto/core"
)

func (tc *Turbine) Process(rs app.Records, fn app.Function) (app.Records, error) {
	return tc.ProcessWithContext(context.Background(), rs, fn)
}

func (tc *Turbine) ProcessWithContext(ctx context.Context, rs app.Records, fn app.Function) (app.Records, error) {
	c, err := tc.AddProcessToCollection(
		ctx,
		&core.ProcessCollectionRequest{
			Process: &core.ProcessCollectionRequest_Process{
				Name: strings.ToLower(reflect.TypeOf(fn).Name()),
			},
			Collection: rs.ToProto(),
		})
	if err != nil {
		return app.Records{}, err
	}

	records := app.NewRecords(c)

	rawOut := fn.Process(records.Records)
	records.Records = rawOut
	return records, nil
}
