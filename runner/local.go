//go:build !platform
// +build !platform

package runner

import (
	"log"

	"github.com/meroxa/turbine-go"

	"github.com/meroxa/turbine-go/local"
)

func Start(app turbine.App) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	lv := local.New()
	err := app.Run(lv)
	if err != nil {
		log.Fatalln(err)
	}
}
