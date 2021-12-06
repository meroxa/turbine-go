package simple

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/meroxa/valve"
	"log"
)

var _ valve.App = (*App)(nil)

type App struct{}

func (a *App) Run(valve valve.Valve) error {
	db, err := valve.Resources("pg")
	if err != nil {
		return err
	}

	rr, err := db.Records("user_activity", nil) // rr is a collection of records, can't be inspected directly
	if err != nil {
		return err
	}

	res, _ := valve.Process(rr, Anonymize{})
	// second return is dead-letter queue

	log.Println("res:", res)

	dwh, err := valve.Resources("sfdwh")
	err = dwh.Write(res, "user_activity", nil)
	if err != nil {
		return err
	}

	return nil
}

type Anonymize struct{}

func (f Anonymize) Process(rr []valve.Record) ([]valve.Record, []valve.RecordWithError) {
	for _, r := range rr {
		r.Payload["email"] = consistentHash(r.Payload["email"].(string))
	}
	return rr, nil
}

func consistentHash(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}
