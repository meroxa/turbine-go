package main

import (
	"github.com/meroxa/valve/examples/simple"
	"github.com/meroxa/valve/local"
)

func main() {
	a := simple.App{}

	lv := local.New("./examples/simple/fixtures")
	err := a.Run(lv)
	if err != nil {
		panic(err)
	}
}
