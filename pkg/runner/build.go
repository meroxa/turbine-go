//go:build builder
// +build builder

package runner

import (
	"context"
	"flag"
	"log"
	"os"
	"path"

	sdk "github.com/meroxa/turbine-go/pkg/turbine"
	"github.com/meroxa/turbine-go/pkg/turbine/build"
)

func execPath() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("unable to locate executable path; error: %s", err)
	}
	return path.Dir(exePath)
}

func Start(app sdk.App) {
	var (
		gitSHA     string
		listenAddr string
		appPath    string
	)

	flag.StringVar(&gitSHA, "gitsha", "", "git commit sha used to reference the code deployed")
	flag.StringVar(&listenAddr, "turbine-core-server", "", "address of the turbine core server")
	flag.StringVar(&appPath, "app-path", "", "path to the turbine application")

	flag.Parse()

	if appPath == "" {
		appPath = execPath()
	}

	if listenAddr == "" {
		log.Println("require -turbine-core-server flag unset")
		flag.PrintDefaults()
		return
	}

	if err := build.Run(
		context.Background(),
		app,
		listenAddr,
		gitSHA,
		appPath,
	); err != nil {
		log.Fatalln(err)
	}
}
