package server

import (
	"context"

	sdk "github.com/meroxa/turbine-go/v3/pkg/turbine"
)

var _ sdk.Destination = (*destination)(nil)

type destination struct{}

func (d *destination) Write(rr sdk.Records) error {
	return d.WriteWithContext(context.Background(), rr)
}

func (d *destination) WriteWithContext(_ context.Context, _ sdk.Records) error {
	return nil
}
