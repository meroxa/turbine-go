//go:build server
// +build server

package runner

import (
	"context"
	"flag"
	"log"
	"os"

	sdk "github.com/meroxa/turbine-go/pkg/turbine"
	"github.com/meroxa/turbine-go/pkg/turbine/server"
)

func Start(app sdk.App) {
	var (
		listenAddr string
		funcName   string
	)

	flag.StringVar(&listenAddr, "serve-addr", os.Getenv("MEROXA_FUNCTION_ADDR"), "listen serve address")
	flag.StringVar(&funcName, "serve", "", "name of function to serve")
	flag.Parse()

	requiredFlag(funcName, "required -serve flag unset")
	requiredFlag(listenAddr, "required -serve-addr or set MEROXA_FUNCTION_ADDR env var")

	if err := server.Run(context.Background(), app, listenAddr, funcName); err != nil {
		log.Fatalln(err)
	}
}
