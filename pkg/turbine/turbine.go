package turbine

import (
	"context"
)

type Turbine interface {
	Source(string, string, ...Option) (Source, error)
	SourceWithContext(context.Context, string, string, ...Option) (Source, error)

	Destination(string, string, ...Option) (Destination, error)
	DestinationWithContext(context.Context, string, string, ...Option) (Destination, error)

	Process(Records, Function) (Records, error)
	ProcessWithContext(context.Context, Records, Function) (Records, error)
}
