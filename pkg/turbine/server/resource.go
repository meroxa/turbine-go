package server

import (
	"context"

	sdk "github.com/meroxa/turbine-go/v2/pkg/turbine"
)

var _ sdk.Resource = (*resource)(nil)

type resource struct{}

func (r *resource) Records(collection string, cfg sdk.ConnectionOptions) (sdk.Records, error) {
	return r.RecordsWithContext(context.Background(), collection, cfg)
}

func (r *resource) RecordsWithContext(ctx context.Context, collection string, cfg sdk.ConnectionOptions) (sdk.Records, error) {
	return sdk.Records{}, nil
}

func (r *resource) Write(rr sdk.Records, collection string) error {
	return r.WriteWithConfigWithContext(context.Background(), rr, collection, sdk.ConnectionOptions{})
}

func (r *resource) WriteWithContext(ctx context.Context, rr sdk.Records, collection string) error {
	return r.WriteWithConfigWithContext(ctx, rr, collection, sdk.ConnectionOptions{})
}

func (r *resource) WriteWithConfig(rr sdk.Records, collection string, cfg sdk.ConnectionOptions) error {
	return r.WriteWithConfigWithContext(context.Background(), rr, collection, cfg)
}

func (r *resource) WriteWithConfigWithContext(ctx context.Context, rr sdk.Records, collection string, cfg sdk.ConnectionOptions) error {
	return nil
}
