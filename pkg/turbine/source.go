package turbine

import (
	"context"
)

type Source interface {
	Read() (Records, error)
	ReadWithContext(context.Context) (Records, error)
}
