package main

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/meroxa/valve"
	"github.com/meroxa/valve/runner"
	"log"
)

func main() {
	runner.Start(App{})
}

var _ valve.App = (*App)(nil)

type App struct{}

func (a App) Run(v valve.Valve) error {
	db, err := v.Resources("demopg")
	if err != nil {
		return err
	}

	rr, err := db.Records("user_activity", nil) // rr is a collection of records, can't be inspected directly
	if err != nil {
		return err
	}

	res, _ := v.Process(rr, Anonymize{})
	// second return is dead-letter queue

	//dwh, err := v.Resources("rdwh")
	//cfg := valve.ResourceConfigs{
	//	{"value.converter.schemas.enable", "false"},
	//}
	//err = dwh.Write(res, "anonymized_user_activity", cfg)
	//if err != nil {
	//	return err
	//}

	s3, err := v.Resources("s3")
	err = s3.Write(res, "data-app-archive", nil)
	if err != nil {
		return err
	}

	return nil
}

type Anonymize struct{}

func (f Anonymize) Process(rr []valve.Record) ([]valve.Record, []valve.RecordWithError) {
	for i, r := range rr {
		hashedEmail := consistentHash(r.Payload.Get("payload.email").(string))
		err := r.Payload.Set("payload.email", hashedEmail)
		if err != nil {
			log.Println("error setting value: ", err)
			break
		}
		rr[i] = r
	}
	return rr, nil
}

func consistentHash(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}
