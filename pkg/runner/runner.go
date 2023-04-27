//go:build !server && !builder
// +build !server,!builder

package runner

import (
	"log"

	sdk "github.com/meroxa/turbine-go/pkg/turbine"
)

func Start(app sdk.App) {
	log.Fatalf("undefined start routine")
}
