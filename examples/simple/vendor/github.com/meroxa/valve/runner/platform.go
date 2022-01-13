//go:build platform
// +build platform

package runner

import (
	"flag"
	"github.com/meroxa/valve"
	"github.com/meroxa/valve/platform"
	"log"
	"os"
	"path"
)

var (
	InvokeFunction string
	ServeFunction  string
	ListFunctions  bool
	BuildImage     bool
	DeployApp      bool
	Help           bool
)

func Start(app valve.App) {

	flag.StringVar(&InvokeFunction, "function", "", "function to trigger")
	flag.StringVar(&ServeFunction, "serve", "", "serve function via gRPC")
	flag.BoolVar(&ListFunctions, "listfunctions", false, "list available functions")
	flag.BoolVar(&BuildImage, "buildimage", false, "build docker image")
	flag.BoolVar(&Help, "help", false, "display help") // TODO: make this trigger by default
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
			log.Fatalf("unable to serve function %s; error: %s", ServeFunction, err)
		}
	}

	if ListFunctions {
		log.Printf("available functions: %s", pv.ListFunctions())
	}

	if BuildImage {
		exePath, err := os.Executable()
		if err != nil {
			log.Fatalf("unable to locate executable path; error: %s", err)
		}

		projPath := path.Dir(exePath)
		projName := path.Base(exePath)
		pv.BuildDockerImage(projName, projPath)
	}

}
