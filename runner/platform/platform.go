package platform

import (
	"github.com/meroxa/turbine-go/runner/unimplemented"
)

type Runner struct {
	unimplemented.Runner
}

func NewRunner() Runner {
	return Runner{}
}
