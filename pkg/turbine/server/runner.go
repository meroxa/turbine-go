package server

import (
	"context"

	sdk "github.com/meroxa/turbine-go/pkg/turbine"
)

func Run(_ context.Context, app sdk.App, addr, fn string) error {
	s := NewServer()
	if err := app.Run(s); err != nil {
		return err
	}

	return s.Listen(addr, fn)
}
