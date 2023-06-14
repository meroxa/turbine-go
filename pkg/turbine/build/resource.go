package build

import (
	"context"

	sdk "github.com/meroxa/turbine-go/v2/pkg/turbine"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"github.com/meroxa/turbine-core/pkg/client"
)

type resource struct {
	*pb.Resource
	client.Client
}

func (r *resource) Records(collection string, cfg sdk.ConnectionOptions) (sdk.Records, error) {
	return r.RecordsWithContext(context.Background(), collection, cfg)
}

func (r *resource) RecordsWithContext(ctx context.Context, collection string, cfg sdk.ConnectionOptions) (sdk.Records, error) {
	c, err := r.ReadCollection(ctx, &pb.ReadCollectionRequest{
		Resource:   r.Resource,
		Collection: collection,
		Configs:    connectionOptions(cfg),
	})
	if err != nil {
		return sdk.Records{}, err
	}

	return collectionToRecords(c), nil
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
	_, err := r.WriteCollectionToResource(ctx, &pb.WriteCollectionRequest{
		Resource:         r.Resource,
		SourceCollection: recordsToCollection(rr),
		TargetCollection: collection,
		Configs:          connectionOptions(cfg),
	})

	return err
}
