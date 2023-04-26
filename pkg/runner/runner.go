package runner

import (
	"context"
	"flag"
	"log"
	"os"
	"path"

	"github.com/meroxa/turbine-go/pkg/turbine"
	"github.com/meroxa/turbine-go/pkg/turbine/build"
)

var (
	gitSha            string
	turbineCoreServer string
	appPath           string
)

func execPath() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("unable to locate executable path; error: %s", err)
	}
	return path.Dir(exePath)
}

func Start(app turbine.App) {
	ctx := context.Background()

	flag.StringVar(&gitSha, "gitsha", "", "git commit sha used to reference the code deployed")
	flag.StringVar(&turbineCoreServer, "turbine-core-server", "", "address of the turbine core server")
	flag.StringVar(&appPath, "app-path", "", "path to the turbine application")
	flag.Parse()

	if appPath == "" {
		appPath = execPath()
	}

	b, err := build.NewBuildClient(ctx, turbineCoreServer, gitSha, appPath)
	if err != nil {
		log.Fatalln(err)
	}

	if err = app.Run(b); err != nil {
		log.Fatalln(err)
	}
}
