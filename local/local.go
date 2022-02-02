package local

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/meroxa/valve"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"time"
	"unsafe"
)

type Valve struct {
	config valve.AppConfig
}

func New() Valve {
	ac, err := valve.ReadAppConfig()
	if err != nil {
		log.Fatalln(err)
	}
	return Valve{ac}
}

func (v Valve) Resources(name string) (valve.Resource, error) {
	return Resource{
		Name:         name,
		fixturesPath: v.config.Resources[name],
	}, nil
}

func (v Valve) Process(rr valve.Records, fn valve.Function) (valve.Records, valve.RecordsWithErrors) {
	var out valve.Records
	var outE valve.RecordsWithErrors

	// use reflection to access intentionally hidden fields
	inVal := reflect.ValueOf(&rr).Elem().FieldByName("records")

	// hack to create reference that can be accessed
	in := reflect.NewAt(inVal.Type(), unsafe.Pointer(inVal.UnsafeAddr())).Elem()
	inRR := in.Interface().([]valve.Record)

	rawOut, _ := fn.Process(inRR)
	out = valve.NewRecords(rawOut)

	return out, outE
}

type Resource struct {
	Name         string
	fixturesPath string
}

func (r Resource) Records(collection string, cfg valve.ResourceConfigs) (valve.Records, error) {
	return readFixtures(r.fixturesPath, collection)
}

func (r Resource) Write(rr valve.Records, collection string, cfg valve.ResourceConfigs) error {
	prettyPrintRecords(r.Name, collection, valve.GetRecords(rr))
	return nil
}

func prettyPrintRecords(name string, collection string, rr []valve.Record) {
	for _, r := range rr {
		log.Printf("%s (%s) => Key: %s; Payload: %s; Timestamp: %s\n", name, collection, r.Key, string(r.Payload), r.Timestamp)
	}
}

type fixtureRecord struct {
	Key       string
	Value     map[string]interface{}
	Timestamp string
}

func readFixtures(path, collection string) (valve.Records, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return valve.Records{}, err
	}

	var records map[string][]fixtureRecord
	err = json.Unmarshal(b, &records)
	if err != nil {
		return valve.Records{}, err
	}

	var rr []valve.Record
	for _, r := range records[collection] {
		rr = append(rr, wrapRecord(r))
	}

	return valve.NewRecords(rr), nil
}

func mapFixturesPath(name, path string) string {
	return fmt.Sprintf("%s/%s.json", path, name)
}

func wrapRecord(m fixtureRecord) valve.Record {
	b, _ := json.Marshal(m.Value)

	var t time.Time
	if m.Timestamp == "" {
		t = time.Now()
	} else {
		// TODO: parse timestamp
	}
	return valve.Record{
		Key:       m.Key,
		Payload:   b,
		Timestamp: t,
	}
}

// Secrets pulls envionment variables with the same name
func (v Valve) RegisterSecret(name string) error {
	val := os.Getenv(name)
	if val == "" {
		return errors.New("secret is invalid or not set")
	}

	return nil
}
