package local

import (
	"context"
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/meroxa/turbine-go"
	"github.com/meroxa/turbine-go/models"
	"log"
	"reflect"

	_ "github.com/apache/beam/sdks/v2/go/pkg/beam/runners/direct"
)

type Runner struct {
	config   turbine.AppConfig
	pipeline *beam.Pipeline
	scope    beam.Scope
}

func (r Runner) RegisterSecret(s string) error {
	//TODO implement me
	panic("implement me")
}

func (r Runner) Secret(s string) string {
	//TODO implement me
	panic("implement me")
}

func NewRunner() Runner {
	ac, err := turbine.ReadAppConfig("", "")
	if err != nil {
		log.Fatalln(err)
	}

	beam.RegisterType(reflect.TypeOf((*models.Record)(nil)).Elem())

	p := beam.NewPipeline()
	s := p.Root()
	return Runner{
		config:   ac,
		pipeline: p,
		scope:    s,
	}
}

func (r Runner) Run() {
	if res, err := beam.Run(context.Background(), "direct", r.pipeline); err != nil {
		log.Fatalf("Failed to execute job: %v", err)
	} else {
		log.Printf("res: %+v", res)
	}
}
