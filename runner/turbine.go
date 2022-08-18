package runner

import (
	"flag"
	"github.com/meroxa/turbine-go"
	"github.com/meroxa/turbine-go/runner/local"
	"log"
)

type App struct {
	turbine.Runner
}

func New() App {
	tr := Init()
	return App{
		tr,
	}
}

var (
	TargetRunnerFlag string
)

const (
	LocalRunner    = "local"
	PlatformRunner = "platform"
	InfoRunner     = "info"
)

// pick the appropriate runner
func Init() turbine.Runner {
	flag.StringVar(&TargetRunnerFlag, "runner", "local", "target runner for execution")
	flag.Parse()

	switch TargetRunnerFlag {
	case LocalRunner:
		return local.NewRunner()
	//case InfoRunner:
	//	TargetRunner = info.NewRunner()
	//case PlatformRunner:
	//	TargetRunner = platform.NewRunner()
	default:
		log.Fatalf("uknown or unsupported runner %s", TargetRunnerFlag)
		return nil
	}
}
