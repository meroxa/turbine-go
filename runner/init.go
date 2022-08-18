package runner

//
//import (
//	"flag"
//	"github.com/meroxa/turbine-go"
//	"github.com/meroxa/turbine-go/runner/local"
//	"log"
//)
//
//var (
//	TargetRunnerFlag string
//
//	TargetRunner turbine.Runner
//)
//
//const (
//	LocalRunner    = "local"
//	PlatformRunner = "platform"
//	InfoRunner     = "info"
//)
//
//// pick the appropriate runner
//func Init() {
//	flag.StringVar(&TargetRunnerFlag, "runner", "local", "target runner for execution")
//	flag.Parse()
//
//	switch TargetRunnerFlag {
//	case LocalRunner:
//		TargetRunner = local.NewRunner()
//	//case InfoRunner:
//	//	TargetRunner = info.NewRunner()
//	//case PlatformRunner:
//	//	TargetRunner = platform.NewRunner()
//	default:
//		log.Fatalf("uknown or unsupported runner %s", TargetRunnerFlag)
//	}
//}
