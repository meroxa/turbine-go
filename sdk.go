package valve

import (
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/valve/internal"
	"log"
)

var Client meroxa.Client

func init() {
	// create global client
	c, err := internal.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	Client = c
}
