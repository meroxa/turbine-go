package v2

import (
	"github.com/gofrs/uuid"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/turbine-go"
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

// TODO: Check if it's still valid
type resourceWithCollection struct {
	Source      bool
	Destination bool
	Name        string
	Collection  string
}

// TODO: Implement
func (t *Turbine) Resources(name string) (turbine.Resource, error) {
	return nil, nil
}

// TODO: Implement
func (t Turbine) ListResources() ([]resourceWithCollection, error) {
	// TODO
	return nil, nil
}

// TODO: Implement
func (r *Resource) Records(collection string, cfg turbine.ResourceConfigs) (turbine.Records, error) {
	return turbine.Records{}, nil
}

// TODO: Implement
func (r *Resource) Write(rr turbine.Records, collection string) error {
	r.Collection = collection
	r.Destination = true
	return r.WriteWithConfig(rr, collection, turbine.ResourceConfigs{})
}

// TODO: Implement
func (r *Resource) WriteWithConfig(rr turbine.Records, collection string, cfg turbine.ResourceConfigs) error {
	return nil
}
