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

	if funcName == "" {
		log.Println("required -serve flag unset")
		flag.PrintDefaults()
		return
	}

	if err := server.Run(context.Background(), app, listenAddr, funcName); err != nil {
		log.Fatalln(err)
	}
}
