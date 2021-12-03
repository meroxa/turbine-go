package platform

import (
	"context"
	"github.com/google/uuid"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/valve"
	"log"
)

type Valve struct {
	client meroxa.Client
}

func New() Valve {
	c, err := newClient()
	if err != nil {
		log.Fatalln(err)
	}
	return Valve{
		client: c,
	}
}

func (v Valve) Resources(name string) (valve.Resource, error) {
	cr, err := v.client.GetResourceByNameOrID(context.Background(), name)
	if err != nil {
		return nil, err
	}

	return Resource {
		ID: cr.ID,
		Name: cr.Name,
		Type: string(cr.Type),
	}, nil
}

type Resource struct {
	ID   int
	UUID uuid.UUID
	Name string
	Type string
}

func (r Resource) Records(collection string, cfg valve.ResourceConfigs) (valve.Records, error) {
	return valve.Records{}, nil
}

func (r Resource) Write(rr valve.Records, collection string, cfg valve.ResourceConfigs) error {
	return nil
}

func (v Valve) Process(rr valve.Records, fn valve.Function) (valve.Records, valve.RecordsWithErrors) {
	var out valve.Records
	var outE valve.RecordsWithErrors

	return out, outE
}
