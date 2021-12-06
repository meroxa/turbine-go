package main

import (
	"github.com/meroxa/valve"
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

	// trigger function fn1
	pv.TriggerFunction("fn1", []valve.Record{})
}