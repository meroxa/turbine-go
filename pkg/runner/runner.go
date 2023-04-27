//go:build !server && !builder
// +build !server,!builder

package runner

import (
	sdk "github.com/meroxa/turbine-go/pkg/turbine"
)

func Start(app sdk.App) {
	requiredFlag("undefined", "undefined start routine")
}
