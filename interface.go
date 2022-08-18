package turbine

import "github.com/meroxa/turbine-go/models"

type Runner interface {
	Resource(string) (Resource, error)
	//Process(Records, Function) Records
	Process(models.Records, models.BeamFunction) models.Records
	RegisterSecret(string) error
	Secret(string) string
	Run()
}

type Resource interface {
	Records(string) (models.Records, error)
	RecordsWithConfig(string, models.ResourceConfig) (models.Records, error)
	Write(models.Records, string) error
	WriteWithConfig(models.Records, string, models.ResourceConfig) error
}
