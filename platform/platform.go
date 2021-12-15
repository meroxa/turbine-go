package platform

import (
	"context"
	"github.com/google/uuid"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/valve"
	"log"
	"reflect"
)

type Valve struct {
	client    meroxa.Client
	functions map[string]valve.Function
}

func New() Valve {
	c, err := newClient()
	if err != nil {
		log.Fatalln(err)
	}
	return Valve{
		client:    c,
		functions: make(map[string]valve.Function),
	}
}

func (v Valve) Resources(name string) (valve.Resource, error) {
	cr, err := v.client.GetResourceByNameOrID(context.Background(), name)
	if err != nil {
		return nil, err
	}

	log.Printf("retrieved resource %s (%s)", cr.Name, cr.Type)

	return Resource{
		ID:     cr.ID,
		Name:   cr.Name,
		Type:   string(cr.Type),
		client: v.client,
	}, nil
}

type Resource struct {
	ID     int
	UUID   uuid.UUID
	Name   string
	Type   string
	client meroxa.Client
}

func (r Resource) Records(collection string, cfg valve.ResourceConfigs) (valve.Records, error) {
	ci := &meroxa.CreateConnectorInput{
		ResourceID:    r.ID,
		Configuration: cfg.ToMap(),
		Type:          meroxa.ConnectorTypeSource,
		Input:         collection,
		PipelineName:  "default",
	}

	con, err := r.client.CreateConnector(context.Background(), ci)
	if err != nil {
		return valve.Records{}, err
	}

	log.Printf("streams: %+v", con.Streams)
	outStreams := con.Streams["output"].([]interface{})

	// Get first output stream
	out := outStreams[0].(string)

	log.Printf("created source connector to resource %s and write records to stream %s to collection %s", r.Name, out, collection)
	return valve.Records{
		Stream: out,
	}, nil
}

func (r Resource) Write(rr valve.Records, collection string, cfg valve.ResourceConfigs) error {
	ci := &meroxa.CreateConnectorInput{
		ResourceID:    r.ID,
		Configuration: cfg.ToMap(),
		Type:          meroxa.ConnectorTypeDestination,
		Input:         rr.Stream,
		PipelineName:  "default",
	}

	// TODO: Apply correct configuration to specify target collection

	_, err := r.client.CreateConnector(context.Background(), ci)
	if err != nil {
		return err
	}
	log.Printf("created destination connector to resource %s and write records from stream %s to collection %s", r.Name, rr.Stream, collection)
	return nil
}

func (v Valve) Process(rr valve.Records, fn valve.Function) (valve.Records, valve.RecordsWithErrors) {
	// TODO: Deploy function
	log.Printf("Deploy function with input stream %s", rr.Stream)
	var out valve.Records
	var outE valve.RecordsWithErrors

	// register function
	v.functions[reflect.TypeOf(fn).Name()] = fn
	out.Stream = uuid.NewString()

	out = rr
	return out, outE
}

func (v Valve) TriggerFunction(name string, in []valve.Record) ([]valve.Record, []valve.RecordWithError) {
	log.Printf("Triggered function %s", name)
	return nil, nil
}
