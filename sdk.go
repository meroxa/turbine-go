package valve

import (
	"github.com/meroxa/valve/internal"
	"log"
)

var Client internal.Client

func init() {
	// create global client
	c, err := internal.NewClient(true)
	if err != nil {
		log.Fatal(err)
	}

	Client = c
}
