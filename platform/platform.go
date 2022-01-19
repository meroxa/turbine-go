package platform

import (
	"context"
	"github.com/google/uuid"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/valve"
	"log"
	"os"
	"path"
	"reflect"
	"strings"
)

type Valve struct {
	client     *Client
	functions  map[string]valve.Function
	deploy     bool
	builtImage string
	pipeline   string
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
	funcName := strings.ToLower(reflect.TypeOf(fn).Name())
	v.functions[funcName] = fn

	var out valve.Records
	var outE valve.RecordsWithErrors

	if v.deploy {
		// ensure that the container image is built and deployed
		if v.builtImage == "" {
			imageName, err := v.buildAndPushFunctionImage()
			log.Printf("building image %s ...", imageName)
			if err != nil {
				log.Panicf("unable to build and push image; err: %s", err.Error())
			}
			v.builtImage = imageName
			log.Printf("image %s build complete", imageName)
		} else {
			log.Printf("image %s already built, using existing image", v.builtImage)
		}

		// create the function
		cfi := CreateFunctionInput{
			InputStream: rr.Stream,
			Image:       v.builtImage,
			EnvVars:     nil,
			Args:        []string{funcName},
			Pipeline:    "default",
		}

		log.Printf("creating function %s ...", funcName)
		fnOut, err := v.client.CreateFunction(context.Background(), &cfi)
		if err != nil {
			log.Printf("create function request: %+v", cfi)
			log.Panicf("unable to build and push image; err: %s", err.Error())
		}
		log.Printf("function %s created (%s)", funcName, fnOut.UUID)
		out.Stream = fnOut.OutputStream
	} else {
		// Not deploying, so map input stream to output stream
		out = rr
	}

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

func (v Valve) buildAndPushFunctionImage() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("unable to locate executable path; error: %s", err)
	}

	projPath := path.Dir(exePath)
	projName := path.Base(exePath)
	v.BuildDockerImage(projName, projPath)

	v.PushDockerImage(projName)
	return projName, nil
}
