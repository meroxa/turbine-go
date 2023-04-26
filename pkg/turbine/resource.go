package turbine

import (
	"context"
)

type Resource interface {
	Records(string, ConnectionOptions) (Records, error)
	RecordsWithContext(context.Context, string, ConnectionOptions) (Records, error)

	Write(Records, string) error
	WriteWithContext(context.Context, Records, string) error
	WriteWithConfig(Records, string, ConnectionOptions) error
	WriteWithConfigWithContext(context.Context, Records, string, ConnectionOptions) error
}
