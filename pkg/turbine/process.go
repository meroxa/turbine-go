package turbine

import (
	"context"
	"reflect"
	"strings"

	"github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
)

func (tc *turbine) Process(rs Records, fn Function) (Records, error) {
	return tc.ProcessWithContext(context.Background(), rs, fn)
}

func (tc *turbine) ProcessWithContext(ctx context.Context, rs Records, fn Function) (Records, error) {
	c, err := tc.AddProcessToCollection(
		ctx,
		&core.ProcessCollectionRequest{
			Process: &core.ProcessCollectionRequest_Process{
				Name: strings.ToLower(reflect.TypeOf(fn).Name()),
			},
			Collection: rs.ToProto(),
		})
	if err != nil {
		return Records{}, err
	}

	out := NewRecords(c)
	rawOut := fn.Process(out.Records)
	out.Records = rawOut
	return out, nil
}
