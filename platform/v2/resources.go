package v2

import (
	"fmt"

	"github.com/meroxa/turbine-go"
	"github.com/meroxa/turbine-go/platform"
)

type Resource struct {
	Name        string
	Source      bool
	Destination bool
	Collection  string
	v           *Turbine
}

func (t *Turbine) Resources(name string) (turbine.Resource, error) {
	r := &Resource{
		Name: name,
		v:    t,
	}
	t.resources = append(t.resources, r)
	return r, nil
}

func (t *Turbine) ListResources() ([]platform.ResourceWithCollection, error) {
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
	records := turbine.Records{}
	if collection == "" {
		return records, fmt.Errorf("please provide a collection name to Records()")
	}

	r.Collection = collection
	r.Source = true

	for _, c := range r.v.deploySpec.Connectors {
		// Only one source per app allowed.
		if c.Type == "source" {
			return records, fmt.Errorf("only one call to Records() is allowed per Meroxa Data Application")
		}
	}

	r.v.deploySpec.Connectors = append(r.v.deploySpec.Connectors,
		specConnector{Type: "source", Resource: r.Name, Collection: collection, Config: cfg.ToMap()})
	return records, nil
}

func (r *Resource) Write(rr turbine.Records, collection string) error {
	if collection == "" {
		return fmt.Errorf("please provide a collection name to Write()")
	}
	return r.WriteWithConfig(rr, collection, turbine.ResourceConfigs{})
}

func (r *Resource) WriteWithConfig(rr turbine.Records, collection string, cfg turbine.ResourceConfigs) error {
	// This function may be called zero or more times.
	if collection == "" {
		return fmt.Errorf("please provide a collection name to WriteWithConfig()")
	}
	r.Collection = collection
	r.Destination = true

	r.v.deploySpec.Connectors = append(r.v.deploySpec.Connectors,
		specConnector{Type: "destination", Resource: r.Name, Collection: collection, Config: cfg.ToMap()})
	return nil
}
