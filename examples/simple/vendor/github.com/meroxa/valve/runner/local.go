//go:build !platform
// +build !platform

package runner

import (
	"github.com/meroxa/valve"
	"github.com/meroxa/valve/local"
	"log"
)

func Start(app valve.App) {
	lv := local.New("./fixtures")
	err := app.Run(lv)
	if err != nil {
		log.Fatalln(err)
	}
}
