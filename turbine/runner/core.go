package runner

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/meroxa/turbine-go/platform"
	"github.com/meroxa/turbine-go/turbine"
	turbineApp "github.com/meroxa/turbine-go/turbine/app"
	"github.com/meroxa/turbine-go/turbine/core"
)

var (
	// Deploy        bool
	// GitSha        string
	// ImageName     string
	// AppName       string
	// ListFunctions bool
	// ListResources bool
	serveFunction string
	run           bool
	record        bool
)

type TurbinePlatformRunner interface {
	turbine.Turbine
	//  run, record, build, serve
	// GetFunction(name string) (turbine.Function, bool)
	// ListFunctions() []string
	// ListResources() ([]platform.ResourceWithCollection, error)
	// DeploymentSpec() (string, error)
}

func Start(app turbine.App) {
	flag.StringVar(&serveFunction, "serve", "", "serve function via gRPC")
	// flag.BoolVar(&ListFunctions, "listfunctions", false, "list available functions")
	// flag.BoolVar(&ListResources, "listresources", false, "list currently used resources")
	// flag.BoolVar(&Deploy, "deploy", false, "deploy the data app")
	// flag.StringVar(&ImageName, "imagename", "", "image name of function image")
	// flag.StringVar(&AppName, "appname", "", "name of application")
	// flag.StringVar(&GitSha, "gitsha", "", "git commit sha used to reference the code deployed")
	flag.BoolVar(&run, "run", false, "run in local mode")
	flag.BoolVar(&record, "record", false, "run in recording mode")
	flag.Parse()

	fmt.Println("run", run)
	fmt.Println("record", record)
	if run {
		err := initializeServer(app, false)
		if err != nil {
			log.Fatalln(err)
		}
	}
	if record {
		err := initializeServer(app, true)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if serveFunction != "" {
		fn, ok := pv.GetFunction(serveFunction)
		if !ok {
			log.Fatalf("invalid or missing function %s", serveFunction)
		}
		err := platform.ServeFunc(fn)
		if err != nil && err.Error() != "received signal terminated" {
			log.Fatalf("unable to serve function %s; error: %s", serveFunction, err)
		}
	}

}

func initializeServer(app turbine.App, recording bool) error {
	pv, err := turbineApp.New(context.TODO(), recording)
	if err != nil {
		return err
	}

	binPath, err := os.Executable()
	if err != nil {
		return err
	}

	pv.Init(context.TODO(), &core.InitRequest{
		AppName:        "test-go",
		ConfigFilePath: path.Dir(binPath),
		Language:       core.Language_GOLANG,
		GitSHA:         "somesha",
		TurbineVersion: "0.1.1",
	})

	return app.Run(pv)
}
