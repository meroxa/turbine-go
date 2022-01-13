package platform

import (
	"context"
	"github.com/google/uuid"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/valve"
	"log"
	"reflect"
	"strings"
)

type Valve struct {
	client    meroxa.Client
	functions map[string]valve.Function
	deploy    bool
}

func New(deploy bool) Valve {
	c, err := newClient()
	if err != nil {
		log.Fatalln(err)
	}
	return Valve{
		client:    c,
		functions: make(map[string]valve.Function),
		deploy:    deploy,
	}
}

func (v Valve) Resources(name string) (valve.Resource, error) {
	if !v.deploy {
		return Resource{}, nil
	}
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
	if r.client == nil {
		return valve.Records{}, nil
	}
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
	if r.client == nil {
		return nil
	}
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
	// register function
	v.functions[strings.ToLower(reflect.TypeOf(fn).Name())] = fn

	if v.deploy {
		// TODO: Deploy function
		log.Printf("TODO: Deploy function with input stream %s", rr.Stream)
	}

	var out valve.Records
	var outE valve.RecordsWithErrors
	out.Stream = uuid.NewString()

	out = rr
	return out, outE
}

func (v Valve) TriggerFunction(name string, in []valve.Record) ([]valve.Record, []valve.RecordWithError) {
	log.Printf("Triggered function %s", name)
	return nil, nil
}

func (v Valve) GetFunction(name string) (valve.Function, bool) {
	fn, ok := v.functions[name]
	return fn, ok
}

func (v Valve) ListFunctions() []string {
	var funcNames []string
	for name := range v.functions {
		funcNames = append(funcNames, name)
	}

	return funcNames
}

func buildAndPushFunctionImage() error {
	return nil
}
