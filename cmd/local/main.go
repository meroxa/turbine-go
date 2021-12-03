package main

import "github.com/meroxa/valve/examples/simple"

func main() {
	a := simple.App{}

	err := a.Run()
	if err != nil {
		panic(err)
	}
}
