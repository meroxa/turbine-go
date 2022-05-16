//go:build platform
// +build platform

package runner

import (
	"flag"
	"log"

	"github.com/meroxa/turbine-go"
	"github.com/meroxa/turbine-go/platform"
)

var (
	ServeFunction string
	ListFunctions bool
	Deploy        bool
	ImageName     string
)

func Start(app turbine.App) {
	flag.StringVar(&ServeFunction, "serve", "", "serve function via gRPC")
	flag.BoolVar(&ListFunctions, "listfunctions", false, "list available functions")
	flag.BoolVar(&Deploy, "deploy", false, "deploy the data app")
	flag.StringVar(&ImageName, "imagename", "", "image name of function image")
	flag.Parse()

	pv := platform.New(Deploy, ImageName)
	err := app.Run(pv)
	if err != nil {
		log.Fatalln(err)
	}

	if ServeFunction != "" {
		fn, ok := pv.GetFunction(ServeFunction)
		if !ok {
			log.Fatalf("invalid or missing function %s", ServeFunction)
		}
		err := platform.ServeFunc(fn)
		if err != nil {
			log.Fatalf("unable to serve function %s; error: %s", ServeFunction, err)
		}
	}

	if ListFunctions {
		log.Printf("available functions: %s", pv.ListFunctions())
	}
}
