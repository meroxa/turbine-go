package main

import (
	"github.com/meroxa/valve"
	"hash/fnv"
	"log"
)

func main() {
	db, err := valve.Resources("pg")
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

	dwh, err := valve.Resources("sfdwh")
	err = dwh.Write(res, "user_activity", nil)
	if err != nil {
		log.Fatal(err)
	}
}

type Anonymize struct{}

func (f Anonymize) Process(rr []valve.Record) ([]valve.Record, error) {
	for _, r := range rr {
		r.Payload["email"] = consistentHash(r.Payload["email"].(string))
	}
	return rr, nil
}

func consistentHash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return string(h.Sum32())
}