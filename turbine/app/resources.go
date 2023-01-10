package app

import (
	"context"
	"encoding/json"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/meroxa/turbine-go/turbine"
	"github.com/meroxa/turbine-go/turbine/core"
)

type Resource struct {
	*core.Resource
	v *Turbine
}

func (t *Turbine) Resources(name string) (turbine.Resource, error) {
	r, err := t.GetResource(context.TODO(), &core.GetResourceRequest{
		Name: name,
	})
	if err != nil {
		return nil, err
	}

	return &Resource{
		Resource: r,
		v:        t,
	}, nil
}

func (r *Resource) Records(collection string, cfg turbine.ConnectionOptions) (turbine.Records, error) {
	c, err := r.v.ReadCollection(context.TODO(), &core.ReadCollectionRequest{
		Resource:   r.Resource,
		Collection: collection,
		Configs:    cfg.ToProto(),
	})
	if err != nil {
		return turbine.Records{}, err
	}

	return turbine.NewRecords(c), nil
}

func (r *Resource) Write(rr turbine.Records, collection string) error {
	return r.WriteWithConfig(rr, collection, turbine.ConnectionOptions{})
}

func (r *Resource) WriteWithConfig(rr turbine.Records, collection string, cfg turbine.ConnectionOptions) error {
	_, err := r.v.WriteCollectionToResource(context.TODO(), &core.WriteCollectionRequest{
		Resource:         r.Resource,
		SourceCollection: rr.ToProto(),
		TargetCollection: collection,
		Configs:          cfg.ToProto(),
	})

	prettyPrintRecords(r.Resource.Name, collection, rr.Records)
	return err
}

func prettyPrintRecords(name string, collection string, rr []turbine.Record) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle("Destination %s/%s", name, collection)
	t.AppendHeader(table.Row{"index", "record"})
	for i, r := range rr {
		payloadVal := string(r.Payload)
		m, err := r.Payload.Map()
		if err == nil {
			b, err := json.MarshalIndent(m, "", "    ")
			if err == nil {
				payloadVal = string(b)
			}
		}
		t.AppendRow(table.Row{i, payloadVal})
	}
	t.AppendFooter(table.Row{"records written", len(rr)})
	t.Render()
}
