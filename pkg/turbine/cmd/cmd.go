package cmd

import (
	"context"
	"flag"
	"log"
	"os"

	sdk "github.com/meroxa/turbine-go/v2/pkg/turbine"
	"github.com/meroxa/turbine-go/v2/pkg/turbine/build"
	"github.com/meroxa/turbine-go/v2/pkg/turbine/server"
)

const (
	serverCmdName = "server"
	buildCmdName  = "build"
)

var (
	buildCmd  = flag.NewFlagSet(serverCmdName, flag.ExitOnError)
	serverCmd = flag.NewFlagSet(buildCmdName, flag.ExitOnError)

	serverListenAddr string
	serverFuncName   string
	buildGitSHA      string
	buildListenAddr  string
	buildAppPath     string
	buildRunProcess  bool
)

func Start(app sdk.App) {
	var (
		ctx = context.Background()
		cmd = parseFlags()
	)
	log.SetOutput(os.Stdout)

	switch cmd {
	case serverCmdName:
		if err := server.Run(
			ctx,
			app,
			serverListenAddr,
			serverFuncName,
		); err != nil {
			log.Fatalln(err)
		}
	case buildCmdName:
		if err := build.Run(
			ctx,
			app,
			buildListenAddr,
			buildGitSHA,
			buildAppPath,
			buildRunProcess,
		); err != nil {
			log.Fatalln(err)
		}
	}
}
