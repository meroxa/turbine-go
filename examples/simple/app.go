package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/meroxa/valve"
	"github.com/meroxa/valve/runner"
	"log"
)

func main() {
	runner.Start(App{})
}

var _ valve.App = (*App)(nil)

type App struct{}

func (a App) Run(valve valve.Valve) error {
	db, err := valve.Resources("demopg")
	if err != nil {
		return err
	}

	rr, err := db.Records("user_activity", nil) // rr is a collection of records, can't be inspected directly
	if err != nil {
		return err
	}

	res, _ := valve.Process(rr, Anonymize{})
	// second return is dead-letter queue

	dwh, err := valve.Resources("rdwh")
	err = dwh.Write(res, "anonymized_user_activity", nil)
	if err != nil {
		return err
	}

	//s3, err := valve.Resources("s3")
	//err = s3.Write(res, "", nil)
	//if err != nil {
	//	return err
	//}

	return nil
}

type Anonymize struct{}

func (f Anonymize) Process(rr []valve.Record) ([]valve.Record, []valve.RecordWithError) {
	for i, r := range rr {
		p, err := JSONToMap(r.Payload)
		if err != nil {
			log.Println("error converting to map: ", err)
			break
		}

		p["email"] = consistentHash(p["email"])
		newP, err := MapToJSON(p)
		if err != nil {
			log.Println("error converting to JSON: ", err)
			break
		}

		r.Payload = newP
		rr[i] = r
	}
	return rr, nil
}

func consistentHash(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

func JSONToMap(b []byte) (map[string]string, error) {
	var m map[string]string
	err := json.Unmarshal(b, &m)
	return m, err
}

func MapToJSON(m map[string]string) ([]byte, error) {
	return json.Marshal(m)
}
