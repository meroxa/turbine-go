package runner

import (
	"context"
	"flag"
	"log"

	"github.com/meroxa/turbine-go/pkg/turbine"
)

var (
	gitSha            string
	turbineCoreServer string
)

func Start(app turbine.App) {
	ctx := context.Background()

	flag.StringVar(&gitSha, "gitsha", "", "git commit sha used to reference the code deployed")
	flag.StringVar(&turbineCoreServer, "turbine-core-server", "", "address of the turbine core server")
	flag.Parse()

	cs, err := turbine.NewCoreServer(ctx, turbineCoreServer, gitSha)
	if err != nil {
		log.Fatalln(err)
	}

	if err = app.Run(cs); err != nil {
		log.Fatalln(err)
	}
}
