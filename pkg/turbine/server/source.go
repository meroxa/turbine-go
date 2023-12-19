package server

import (
	"context"

	sdk "github.com/meroxa/turbine-go/v3/pkg/turbine"
)

var _ sdk.Source = (*source)(nil)

type source struct{}

func (r *source) Read() (sdk.Records, error) {
	return r.ReadWithContext(context.Background())
}

func (r *source) ReadWithContext(ctx context.Context) (sdk.Records, error) {
	return sdk.Records{}, nil
}
