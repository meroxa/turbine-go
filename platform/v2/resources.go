package v2

import (
	"fmt"

	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/turbine-go"
	"github.com/meroxa/turbine-go/platform"
)

type Resource struct {
	Name        string
	Source      bool
	Destination bool
	Collection  string
	Connectors  []turbine.SpecConnector
	client      meroxa.Client
	v           *Turbine
}

func (t *Turbine) Resources(name string) (turbine.Resource, error) {
	r := &Resource{
		Name: name,
	}
	t.resources = append(t.resources, r)
	return r, nil
}

func (t Turbine) ListResources() ([]platform.ResourceWithCollection, error) {
	var resources []platform.ResourceWithCollection

	for i := range t.resources {
		r, ok := (t.resources[i]).(*Resource)
		if !ok {
			return nil, fmt.Errorf("Bad resource type.")
		}
		resources = append(resources, platform.ResourceWithCollection{
			Source:      r.Source,
			Destination: r.Destination,
			Collection:  r.Collection,
			Name:        r.Name,
		})

	}
	return resources, nil
}

func (r *Resource) Records(collection string, cfg turbine.ResourceConfigs) (turbine.Records, error) {
	// This function will only be called once because there is only ever one source.
	r.Collection = collection
	r.Source = true

	r.Connectors = append(
		r.Connectors,
		turbine.SpecConnector{Type: "source", Resource: r.Name, Collection: collection, Config: cfg.ToMap()})
	return turbine.Records{}, nil
}

func (r *Resource) Write(rr turbine.Records, collection string) error {
	return r.WriteWithConfig(rr, collection, turbine.ResourceConfigs{})
}

func (r *Resource) WriteWithConfig(rr turbine.Records, collection string, cfg turbine.ResourceConfigs) error {
	// This function may be called zero to many times.
	r.Collection = collection
	r.Destination = true
	r.Connectors = append(
		r.Connectors,
		turbine.SpecConnector{Type: "destination", Resource: r.Name, Collection: collection, Config: cfg.ToMap()})
	return nil
}

func (r *Resource) GetSpecConnectors() []turbine.SpecConnector {
	return r.Connectors
}
