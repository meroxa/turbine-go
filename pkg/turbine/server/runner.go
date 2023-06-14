package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	sdk "github.com/meroxa/turbine-go/v2/pkg/turbine"
)

func Run(_ context.Context, app sdk.App, addr, fn string) error {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	s := NewServer()
	if err := app.Run(s); err != nil {
		return err
	}

	go func() {
		sig := <-sigchan
		log.Printf("received signal %q, exiting", sig)
		s.GracefulStop()
	}()

	log.Printf("starting function %q on address %q", fn, addr)
	return s.Listen(addr, fn)
}
