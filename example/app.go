package main

import (
	"context"
	"github.com/meroxa/turbine-go/models"
	"github.com/meroxa/turbine-go/runner"
	"log"
)

func main() {
	app := runner.New()

	db, err := app.Resource("demopg")
	if err != nil {
		log.Fatalf("error retrieving resource; err: %s", err.Error())
	}

	rr, err := db.Records("user_activity")
	if err != nil {
		log.Fatalf("error retrieving records; err: %s", err.Error())
	}

	res := app.Process(rr, MyFunc{})

	err = db.Write(res, "user_activity_transformed")
	if err != nil {
		log.Fatalf("error retrieving records; err: %s", err.Error())
	}

	app.Run()
}

type MyFunc struct{}

func (f MyFunc) ProcessElement(ctx context.Context, s string, emit func(string)) {
	emit(s)
}

func (f MyFunc) Process(r models.Record) models.Record {
	return models.Record{}
}
