package build

import (
	"context"

	sdk "github.com/meroxa/turbine-go/pkg/turbine"
)

type Runner struct {
	listenAddr string
	gitSha     string
	appPath    string
}

func Run(ctx context.Context, app sdk.App, addr, gitsha, path string) error {
	b, err := NewBuildClient(ctx, addr, gitsha, path)
	if err != nil {
		return err
	}
	return app.Run(b)
}
