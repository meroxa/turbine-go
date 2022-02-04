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

	err = v.RegisterSecret("CLEARBIT_API_KEY") // makes env var available to data app
	if err != nil {
		return err
	}
	res, _ := v.Process(rr, EnrichUserData{})

	//s3, err := v.Resources("s3")
	//err = s3.Write(res, "user_activity_enriched", nil)
	//if err != nil {
	//	return err
	//}

	err = db.Write(res, "user_activity_enriched", nil)
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

type EnrichUserData struct{}

func (f EnrichUserData) Process(rr []valve.Record) ([]valve.Record, []valve.RecordWithError) {
	for i, r := range rr {
		log.Printf("Got email: %s", r.Payload.Get("email"))
		UserDetails, err := EnrichUserEmail(r.Payload.Get("email").(string))
		if err != nil {
			log.Println("error enriching user data: ", err)
			break
		}
		log.Printf("Got UserDetails: %+v", UserDetails)
		err = r.Payload.Set("full_name", UserDetails.FullName)
		err = r.Payload.Set("company", UserDetails.Company)
		err = r.Payload.Set("location", UserDetails.Location)
		err = r.Payload.Set("role", UserDetails.Role)
		err = r.Payload.Set("seniority", UserDetails.Seniority)
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
