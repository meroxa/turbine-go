package runtime

import (
	"context"
)

type Runtime interface {
	FindOrCreatePipeline(ctx context.Context, pipelineName string) (Pipeline, error)
	FindOrCreateApp(ctx context.Context, app App) (App, error)
	FindResource(ctx context.Context, resourceName string) (Resource, error)
	FindOrCreateProcess(ctx context.Context, fn Processor) (Processor, error)
	FindOrCreateInputConnector(ctx context.Context, resource Resource, collectionName string, cfg ResourceConfigs)
	FindOrCreateOutputConnector(ctx context.Context, resource Resource, collectionName string, cfg ResourceConfigs)

	GetProcessor(name string) (Processor, bool)
	ListProcessors() []string
	ListReferencedResources() ([]Resource, error)
	RegisterSecret(name string) error
}
