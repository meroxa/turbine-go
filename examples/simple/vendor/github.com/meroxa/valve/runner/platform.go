//go:build platform
// +build platform

package runner

import (
	"flag"
	"github.com/meroxa/valve"
	"github.com/meroxa/valve/platform"
	"log"
)

var (
	InvokeFunction string
	ServeFunction  string
	ListFunctions  bool
	DeployApp      bool
	Help           bool
)

func Start(app valve.App) {

	flag.StringVar(&InvokeFunction, "functions", "", "function to trigger")
	flag.StringVar(&ServeFunction, "serve", "", "serve function via gRPC")
	flag.BoolVar(&ListFunctions, "listfunctions", false, "list available functions")
	flag.BoolVar(&Help, "help", false, "display help")
	flag.BoolVar(&DeployApp, "deploy", false, "deploy the data app")
	flag.Parse()

	pv := platform.New(DeployApp)
	err := app.Run(pv)
	if err != nil {
		log.Fatalln(err)
	}

	if InvokeFunction != "" {
		pv.TriggerFunction(InvokeFunction, nil)
	}

	if ServeFunction != "" {
		fn, ok := pv.GetFunction(ServeFunction)
		if !ok {
			log.Fatalf("invalid or missing function %s", ServeFunction)
		}
		err := platform.ServeFunc(fn)
		if err != nil {
			log.Fatalf("unable to serve function %s; error: ", ServeFunction, err)
		}
	}

	if ListFunctions {
		log.Printf("available functions: %s", pv.ListFunctions())
	}

}
