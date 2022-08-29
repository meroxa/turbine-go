package v2

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/turbine-go"
	"github.com/meroxa/turbine-go/platform"
)

// Define my own version of Resource
type Resource struct {
	UUID        uuid.UUID
	Name        string
	Type        string
	Source      bool
	Destination bool
	Collection  string
	client      meroxa.Client
	v           *Turbine
}

func (t *Turbine) Resources(name string) (turbine.Resource, error) {
	fmt.Println("will create a connector based on the resource. first time a source, following times destinations")
	return nil, nil
}

// TODO: Implement
func (t Turbine) ListResources() ([]platform.ResourceWithCollection, error) {
	// TODO
	fmt.Println("Will list resources")
	return nil, nil
}

// TODO: Implement
func (r *Resource) Records(collection string, cfg turbine.ResourceConfigs) (turbine.Records, error) {
	fmt.Println("Will attach stream")
	return turbine.Records{}, nil
}

// TODO: Implement
func (r *Resource) Write(rr turbine.Records, collection string) error {
	fmt.Println("Add a destination / writing to prior stream")
	r.Collection = collection
	r.Destination = true
	return r.WriteWithConfig(rr, collection, turbine.ResourceConfigs{})
}

// TODO: Implement
func (r *Resource) WriteWithConfig(rr turbine.Records, collection string, cfg turbine.ResourceConfigs) error {
	fmt.Println("Write with config")
	return nil
}
