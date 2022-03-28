package main

import (
	"github.com/meroxa/turbine"
	"github.com/meroxa/turbine/runner"
)

func main() {
	runner.Start(App{})
}

var _ turbine.App = (*App)(nil)

type App struct{}

func (a App) Run(v turbine.Turbine) error {
    // The Resources function identifies the upstream data store
    // that you want to work with.
	source, err := v.Resources("source_name")
	if err != nil {
		return err
	}

    // The `Records` function will pull
    //
	rr, err := source.Records("collection_name", nil)
	if err != nil {
		return err
	}

	res, _ := v.Process(rr, Anonymize{})

	dest, err := v.Resources("dest_name")
	if err != nil {
		return err
	}

	err = dest.Write(res, "collection_name", nil)
	if err != nil {
		return err
	}

	return nil
}

type Anonymize struct{}

func (f Anonymize) Process(rr []turbine.Record) ([]turbine.Record, []turbine.RecordWithError) {
	return rr, nil
}
