package turbine

import (
	"context"

	"github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"github.com/meroxa/turbine-core/pkg/client"
)

type Resource interface {
	Records(string, ConnectionOptions) (Records, error)
	RecordsWithContext(context.Context, string, ConnectionOptions) (Records, error)

	Write(Records, string) error
	WriteWithContext(context.Context, Records, string) error
	WriteWithConfig(Records, string, ConnectionOptions) error
	WriteWithConfigWithContext(context.Context, Records, string, ConnectionOptions) error
}

type resource struct {
	*core.Resource
	tc client.Client
}

func (tc *turbine) Resources(name string) (Resource, error) {
	return tc.ResourcesWithContext(context.Background(), name)
}

func (tc *turbine) ResourcesWithContext(ctx context.Context, name string) (Resource, error) {
	r, err := tc.GetResource(
		ctx,
		&core.GetResourceRequest{
			Name: name,
		})
	if err != nil {
		return nil, err
	}

	return &resource{
		Resource: r,
		tc:       tc.Client,
	}, nil
}

func (r *resource) Records(collection string, cfg ConnectionOptions) (Records, error) {
	return r.RecordsWithContext(context.Background(), collection, cfg)
}

func (r *resource) RecordsWithContext(ctx context.Context, collection string, cfg ConnectionOptions) (Records, error) {
	c, err := r.tc.ReadCollection(ctx, &core.ReadCollectionRequest{
		Resource:   r.Resource,
		Collection: collection,
		Configs:    cfg.ToProto(),
	})
	if err != nil {
		return Records{}, err
	}

	return NewRecords(c), nil
}

func (r *resource) Write(rr Records, collection string) error {
	return r.WriteWithConfigWithContext(context.Background(), rr, collection, ConnectionOptions{})
}

func (r *resource) WriteWithContext(ctx context.Context, rr Records, collection string) error {
	return r.WriteWithConfigWithContext(ctx, rr, collection, ConnectionOptions{})
}

func (r *resource) WriteWithConfig(rr Records, collection string, cfg ConnectionOptions) error {
	return r.WriteWithConfigWithContext(context.Background(), rr, collection, cfg)
}

func (r *resource) WriteWithConfigWithContext(ctx context.Context, rr Records, collection string, cfg ConnectionOptions) error {
	_, err := r.tc.WriteCollectionToResource(ctx, &core.WriteCollectionRequest{
		Resource:         r.Resource,
		SourceCollection: rr.ToProto(),
		TargetCollection: collection,
		Configs:          cfg.ToProto(),
	})

	return err
}
