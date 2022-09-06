package platform

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"

	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/turbine-go"
)

type Turbine struct {
	client    *Client
	functions map[string]turbine.Function
	resources []turbine.Resource
	deploy    bool
	imageName string
	config    turbine.AppConfig
	secrets   map[string]string
	gitSha    string
	appUUID   string
}

var pipelineUUID string

func New(deploy bool, imageName, appName, gitSha string) *Turbine {
	c, err := NewClient()
	if err != nil {
		log.Fatalln(err)
	}

	ac, err := turbine.ReadAppConfig(appName, "")
	if err != nil {
		log.Fatalln(err)
	}
	return &Turbine{
		client:    c,
		functions: make(map[string]turbine.Function),
		resources: []turbine.Resource{},
		imageName: imageName,
		deploy:    deploy,
		config:    ac,
		secrets:   make(map[string]string),
		gitSha:    gitSha,
	}
}

// TODO: Remove once everything is under IR
func (t *Turbine) findPipeline(ctx context.Context) error {
	_, err := t.client.GetPipelineByName(ctx, t.config.Pipeline)
	return err
}

// TODO: Remove once everything is under IR
func (t *Turbine) createPipeline(ctx context.Context) error {
	input := &meroxa.CreatePipelineInput{
		Name: t.config.Pipeline,
		Metadata: map[string]interface{}{
			"app":     t.config.Name,
			"turbine": true,
		},
	}

	p, err := t.client.CreatePipeline(ctx, input)
	if err != nil {
		return err
	}
	pipelineUUID = p.UUID
	return nil
}

// TODO: Remove once everything is under IR
func (t *Turbine) createApplication(ctx context.Context) error {
	inputCreateApp := &meroxa.CreateApplicationInput{
		Name:     t.config.Name,
		Language: "golang",
		GitSha:   t.gitSha,
		Pipeline: meroxa.EntityIdentifier{UUID: null.StringFrom(pipelineUUID)},
	}
	a, err := t.client.CreateApplication(ctx, inputCreateApp)
	t.appUUID = a.UUID
	return err
}

func (t *Turbine) Resources(name string) (turbine.Resource, error) {
	if !t.deploy {
		r := &Resource{
			Name: name,
		}
		t.resources = append(t.resources, r)
		return r, nil
	}

	ctx := context.Background()

	// Make sure we only create pipeline once
	if ok := t.findPipeline(ctx); ok != nil {
		err := t.createPipeline(ctx)
		if err != nil {
			return nil, err
		}
	}

	resource, err := t.client.GetResourceByNameOrID(ctx, name)
	if err != nil {
		return nil, err
	}

	log.Printf("retrieved resource %s (%s)", resource.Name, resource.Type)

	u, _ := uuid.Parse(resource.UUID)
	return &Resource{
		UUID:   u,
		Name:   resource.Name,
		Type:   string(resource.Type),
		client: t.client,
		v:      t,
	}, nil
}

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

func (r *Resource) Records(collection string, cfg turbine.ResourceConfigs) (turbine.Records, error) {
	r.Collection = collection
	r.Source = true

	if r.client == nil {
		return turbine.Records{}, nil
	}

	connectorConfig := cfg.ToMap()
	switch r.Type {
	case "kafka":
		connectorConfig["conduit"] = "true" // only support Kafka connectors using Conduit so this is safe
	}

	ci := &meroxa.CreateConnectorInput{
		ResourceName:  r.Name,
		Configuration: connectorConfig,
		Type:          meroxa.ConnectorTypeSource,
		Input:         collection,
		PipelineName:  r.v.config.Pipeline,
	}

	con, err := r.client.CreateConnector(context.Background(), ci)
	if err != nil {
		return turbine.Records{}, err
	}

	outStreams := con.Streams["output"].([]interface{})

	// Get first output stream
	out := outStreams[0].(string)

	log.Printf("created source connector to resource %s and write records to stream %s from collection %s", r.Name, out, collection)
	return turbine.Records{
		Stream: out,
	}, nil
}

func (r *Resource) Write(rr turbine.Records, collection string) error {
	r.Collection = collection
	r.Destination = true
	return r.WriteWithConfig(rr, collection, turbine.ResourceConfigs{})
}

func (r *Resource) WriteWithConfig(rr turbine.Records, collection string, cfg turbine.ResourceConfigs) error {
	// bail if dryrun
	if r.client == nil {
		return nil
	}

	connectorConfig := cfg.ToMap()
	switch r.Type {
	case "kafka":
		connectorConfig["conduit"] = "true" // only support Kafka connectors using Conduit so this is safe
		connectorConfig["topic"] = strings.ToLower(collection)
	case "redshift", "postgres", "mysql", "sqlserver": // JDBC sink
		connectorConfig["table.name.format"] = strings.ToLower(collection)
	case "mongodb":
		connectorConfig["collection"] = strings.ToLower(collection)
	case "s3":
		connectorConfig["aws_s3_prefix"] = strings.ToLower(collection) + "/"
	case "snowflakedb":
		r := regexp.MustCompile("^[a-zA-Z]{1}[a-zA-Z0-9_]*$")
		matched := r.MatchString(collection)
		if !matched {
			return fmt.Errorf("%q is an invalid Snowflake name - must start with "+
				"a letter and contain only letters, numbers, and underscores", collection)
		}
		connectorConfig["snowflake.topic2table.map"] =
			fmt.Sprintf("%s:%s", rr.Stream, collection)
	}

	ci := &meroxa.CreateConnectorInput{
		ResourceName:  r.Name,
		Configuration: connectorConfig,
		Type:          meroxa.ConnectorTypeDestination,
		Input:         rr.Stream,
		PipelineName:  r.v.config.Pipeline,
	}

	_, err := r.client.CreateConnector(context.Background(), ci)
	if err != nil {
		return err
	}
	log.Printf("created destination connector to resource %s and write records from stream %s to collection %s", r.Name, rr.Stream, collection)

	if r.v.appUUID == "" {
		err = r.v.createApplication(context.Background())
		if err != nil {
			return err
		}
		log.Printf("created application %q", r.v.config.Name)
	}

	return nil
}

func (t Turbine) Process(rr turbine.Records, fn turbine.Function) turbine.Records {
	// register function and associate it with the last gitsha
	var (
		funcName       = strings.ToLower(reflect.TypeOf(fn).Name())
		funcNameGitSHA = fmt.Sprintf("%s-%.8s", funcName, t.gitSha)
	)

	t.functions[funcName] = fn

	var out turbine.Records

	if t.deploy {
		// create the function
		cfi := &meroxa.CreateFunctionInput{
			Name:        funcNameGitSHA,
			InputStream: rr.Stream,
			Image:       t.imageName,
			EnvVars:     t.secrets,
			Args:        []string{funcName},
			Pipeline:    meroxa.PipelineIdentifier{Name: t.config.Pipeline},
		}

		log.Printf("creating function %s ...", funcName)
		fnOut, err := t.client.CreateFunction(context.Background(), cfi)
		if err != nil {
			log.Panicf("unable to create function; err: %s", err.Error())
		}
		log.Printf("function %s created (%s)", funcName, fnOut.UUID)
		out.Stream = fnOut.OutputStream
	} else {
		// Not deploying, so map input stream to output stream
		out = rr
	}

	return out
}

func (t Turbine) GetFunction(name string) (turbine.Function, bool) {
	fn, ok := t.functions[name]
	return fn, ok
}

func (t Turbine) ListFunctions() []string {
	var funcNames []string
	for name := range t.functions {
		funcNames = append(funcNames, name)
	}

	return funcNames
}

type ResourceWithCollection struct {
	Source      bool
	Destination bool
	Name        string
	Collection  string
}

func (t Turbine) ListResources() ([]ResourceWithCollection, error) {
	var resources []ResourceWithCollection

	for i := range t.resources {
		r, ok := (t.resources[i]).(*Resource)
		if !ok {
			return nil, fmt.Errorf("Bad resource type.")
		}
		resources = append(resources, ResourceWithCollection{
			Source:      r.Source,
			Destination: r.Destination,
			Collection:  r.Collection,
			Name:        r.Name,
		})

	}
	return resources, nil
}

// RegisterSecret pulls environment variables with the same name and ships them as Env Vars for functions
func (t Turbine) RegisterSecret(name string) error {
	val := os.Getenv(name)
	if val == "" {
		return errors.New("secret is invalid or not set")
	}

	t.secrets[name] = val
	return nil
}

func (t Turbine) HandleSpec() (string, error) {
	panic("unimplemented")
}
