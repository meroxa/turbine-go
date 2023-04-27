package runner

import (
	"flag"
	"log"
	"os"
)

func requiredFlag(f, msg string) {
	if f == "" {
		log.Println(msg)
		flag.PrintDefaults()
		os.Exit(1)
	}
}
