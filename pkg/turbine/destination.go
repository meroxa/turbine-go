package turbine

import (
	"context"
)

type Destination interface {
	Write(Records) error
	WriteWithContext(context.Context, Records) error
}
