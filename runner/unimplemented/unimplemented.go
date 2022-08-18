package unimplemented

import (
	"github.com/meroxa/turbine-go/models"
	"github.com/meroxa/turbine-go/runtime"
)

type Runner struct{}

func (u Runner) Resource(resourceName string) (runtime.Resource, error) {
	//TODO implement me
	panic("implement me")
}

func (u Runner) Process(in models.Records, fn models.Function) models.Records {
	//TODO implement me
	panic("implement me")
}

func (u Runner) RegisterSecret(envVar string) error {
	//TODO implement me
	panic("implement me")
}

func (u Runner) Secret(name string) string {
	//TODO implement me
	panic("implement me")
}
