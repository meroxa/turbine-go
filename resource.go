package valve

import "github.com/google/uuid"

type Resource struct {
	Name string
	UUID uuid.UUID
	records []Record
}

type ResourceConfig struct {
	Field string
	Value string
}

type ResourceConfigs []ResourceConfig

func (r Resource) Records(collection string, cfg ResourceConfigs) ([]Record, error) {
	return r.records, nil
}

func (r Resource) Write(rr []Record, collection string, cfg ResourceConfigs) error {
	return nil
}

func Resources(name string) (Resource, error) {
	return Resource{}, nil
}