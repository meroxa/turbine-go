package main

import (
	"github.com/meroxa/valve/examples/simple"
	"github.com/meroxa/valve/platform"
)

func main() {
	a := simple.App{}

	pv := platform.New()
	err := a.Run(pv)
	if err != nil {
		panic(err)
	}

	// trigger function Anonymize
	//pv.TriggerFunction("Anonymize", []valve.Record{})
}