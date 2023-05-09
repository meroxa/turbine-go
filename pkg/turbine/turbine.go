package turbine

import (
	"context"
)

type Turbine interface {
	Resources(string) (Resource, error)
	ResourcesWithContext(context.Context, string) (Resource, error)

	Process(Records, Function) (Records, error)
	ProcessWithContext(context.Context, Records, Function) (Records, error)

	RegisterSecret(name string) error
	RegisterSecretWithContext(ctx context.Context, name string) error
}
