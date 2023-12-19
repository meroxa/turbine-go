package build

import (
	"context"

	sdk "github.com/meroxa/turbine-go/v3/pkg/turbine"
)

func Run(ctx context.Context, app sdk.App, addr, gitsha, path string, runProcess bool) error {
	b, err := NewBuildClient(ctx, addr, gitsha, path, runProcess)
	if err != nil {
		return err
	}
	return app.Run(b)
}
