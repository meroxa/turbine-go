package turbine

import (
	"context"
)

type Turbine interface {
	Server

	RegisterSecret(name string) error
	RegisterSecretWithContext(ctx context.Context, name string) error
}

type Server interface {
	Resources(string) (Resource, error)
	ResourcesWithContext(context.Context, string) (Resource, error)

	Process(Records, Function) (Records, error)
	ProcessWithContext(context.Context, Records, Function) (Records, error)
}
