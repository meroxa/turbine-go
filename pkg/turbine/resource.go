package turbine

import (
	"context"

	"github.com/meroxa/turbine-go/pkg/app"
	"github.com/meroxa/turbine-go/pkg/proto/core"
)

type Resource interface {
	Records(string, app.ConnectionOptions) (app.Records, error)
	RecordsWithContext(context.Context, string, app.ConnectionOptions) (app.Records, error)

	Write(app.Records, string) error
	WriteWithContext(context.Context, app.Records, string) error
	WriteWithConfig(app.Records, string, app.ConnectionOptions) error
	WriteWithConfigWithContext(context.Context, app.Records, string, app.ConnectionOptions) error
}

type resource struct {
	*core.Resource
	tc TurbineCore
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
		tc:       tc.TurbineCore,
	}, nil
}

func (r *resource) Records(collection string, cfg app.ConnectionOptions) (app.Records, error) {
	return r.RecordsWithContext(context.Background(), collection, cfg)
}

func (r *resource) RecordsWithContext(ctx context.Context, collection string, cfg app.ConnectionOptions) (app.Records, error) {
	c, err := r.tc.ReadCollection(ctx, &core.ReadCollectionRequest{
		Resource:   r.Resource,
		Collection: collection,
		Configs:    cfg.ToProto(),
	})
	if err != nil {
		return app.Records{}, err
	}

	return app.NewRecords(c), nil
}

func (r *resource) Write(rr app.Records, collection string) error {
	return r.WriteWithConfigWithContext(context.Background(), rr, collection, app.ConnectionOptions{})
}

func (r *resource) WriteWithContext(ctx context.Context, rr app.Records, collection string) error {
	return r.WriteWithConfigWithContext(ctx, rr, collection, app.ConnectionOptions{})
}

func (r *resource) WriteWithConfig(rr app.Records, collection string, cfg app.ConnectionOptions) error {
	return r.WriteWithConfigWithContext(context.Background(), rr, collection, cfg)
}

func (r *resource) WriteWithConfigWithContext(ctx context.Context, rr app.Records, collection string, cfg app.ConnectionOptions) error {
	_, err := r.tc.WriteCollectionToResource(ctx, &core.WriteCollectionRequest{
		Resource:         r.Resource,
		SourceCollection: rr.ToProto(),
		TargetCollection: collection,
		Configs:          cfg.ToProto(),
	})

	return err
}
