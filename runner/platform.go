//go:build platform
// +build platform

package runner

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/meroxa/turbine-go"
	"github.com/meroxa/turbine-go/platform"
)

var (
	Deploy        bool
	GitSha        string
	ImageName     string
	AppName       string
	ListFunctions bool
	ListResources bool
	ServeFunction string
)

func Start(app turbine.App) {
	flag.StringVar(&ServeFunction, "serve", "", "serve function via gRPC")
	flag.BoolVar(&ListFunctions, "listfunctions", false, "list available functions")
	flag.BoolVar(&ListResources, "listresources", false, "list currently used resources")
	flag.BoolVar(&Deploy, "deploy", false, "deploy the data app")
	flag.StringVar(&ImageName, "imagename", "", "image name of function image")
	flag.StringVar(&AppName, "appname", "", "name of application")
	flag.StringVar(&GitSha, "gitsha", "", "git commit sha used to reference the code deployed")
	flag.Parse()

	pv := platform.New(Deploy, ImageName, AppName, GitSha)

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

	if ListResources {
		rr, err := pv.ListResources()
		if err != nil {
			log.Fatal(err)
		}

		enc := json.NewEncoder(os.Stdout)
		if err := enc.Encode(rr); err != nil {
			log.Fatal(err)
		}
	}
}
