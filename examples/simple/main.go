package main

import (
	"github.com/meroxa/valve"
	"log"
)

func main() {
	db, err := valve.Resources("mypg")
	if err != nil {
		log.Fatal(err)
	}

	rr, err := db.Records("user_activity", nil)
	if err != nil {
		log.Fatal(err)
	}

	res, dl := valve.Process(rr, Anonymize{})
	if len(dl) > 0 { // dead-letter queue not empty
		log.Printf("Error processing %d records", len(dl))
	}

	dwh, err := valve.Resources("dwh")
	err = dwh.Write(res, "user_activity", nil)
	if err != nil {
		log.Fatal(err)
	}
}

type Anonymize struct{}

func (f Anonymize) Process(rr []valve.Record) ([]valve.Record, error) {
	return nil, nil
}
