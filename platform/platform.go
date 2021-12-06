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
	functions map[string]valve.Function
}

func New() Valve {
	c, err := newClient()
	if err != nil {
		log.Fatalln(err)
	}
	return Valve{
		client: c,
		functions: make(map[string]valve.Function),
	}
}

func (v Valve) Resources(name string) (valve.Resource, error) {
	cr, err := v.client.GetResourceByNameOrID(context.Background(), name)
	if err != nil {
		return nil, err
	}

	log.Printf("retrieved resource %s (%s)", cr.Name, cr.Type)

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
	// TODO:
	// - Create source connector
	// - Return valve.Records with output stream set
	return valve.Records{
		Stream: uuid.NewString(),
	}, nil
}

func (r Resource) Write(rr valve.Records, collection string, cfg valve.ResourceConfigs) error {
	// TODO: Create destination connector
	log.Printf("create destination connector to resource %s and write records from stream %s to collection %s", r.Name, rr.Stream, collection)
	return nil
}

func (v Valve) Process(rr valve.Records, fn valve.Function) (valve.Records, valve.RecordsWithErrors) {
	// TODO
	log.Printf("Deploy function with input stream %s", rr.Stream)
	var out valve.Records
	var outE valve.RecordsWithErrors

	// register function
	v.functions["fn1"] = fn

	out.Stream = uuid.NewString()

	return out, outE
}

func (v Valve) TriggerFunction(name string, in []valve.Record) ([]valve.Record, []valve.RecordWithError) {
	log.Printf("Triggered function %s", name)
	return nil, nil
}