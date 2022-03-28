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
	// Identify an upstream data store for your data app
	// with the `Resources` function
	source, err := v.Resources("source_name")
	if err != nil {
		return err
	}

	// Specify which upstream records to pull
	// with the `Records` function
	rr, err := source.Records("collection_name", nil)
	if err != nil {
		return err
	}

	// Specify what code to execute against upstream records
	// with the `Process` function
	res, _ := v.Process(rr, Anonymize{})

	// Identify a downstream data store for your data app
	// with the `Resources` function
	dest, err := v.Resources("dest_name")
	if err != nil {
		return err
	}

	// Specify where to write records downstream
	// using the `Write` function
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
