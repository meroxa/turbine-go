package v2

import (
	"log"

	"github.com/meroxa/turbine-go"
	"github.com/meroxa/turbine-go/platform"
)

//type Turbine interface {
//	Resources(string) (Resource, error)
//	Process(Records, Function) Records
//	RegisterSecret(string) error
//}
//
//type Resource interface {
//	Records(collection string, cfg ResourceConfigs) (Records, error)
//	Write(records Records, collection string) error
//	WriteWithConfig(records Records, collection string, cfg ResourceConfigs) error
//}

// Implement those same methods

// Define my own version of Turbine
type Turbine struct {
	client     *platform.Client
	functions  map[string]turbine.Function
	resources  []turbine.Resource
	deploy     bool
	deploySpec string
	imageName  string
	config     turbine.AppConfig
	secrets    map[string]string
	gitSha     string
}

func New(deploy bool, imageName, appName, gitSha, spec string) *Turbine {
	c, err := platform.NewClient()
	if err != nil {
		log.Fatalln(err)
	}

	ac, err := turbine.ReadAppConfig(appName, "")
	if err != nil {
		log.Fatalln(err)
	}
	return &Turbine{
		client:     c,
		functions:  make(map[string]turbine.Function),
		resources:  []turbine.Resource{},
		imageName:  imageName,
		deploy:     deploy,
		deploySpec: spec,
		config:     ac,
		secrets:    make(map[string]string),
		gitSha:     gitSha,
	}
}

// TODO: Implement
func (t Turbine) Process(rr turbine.Records, fn turbine.Function) turbine.Records {
	var out turbine.Records
	return out
}
