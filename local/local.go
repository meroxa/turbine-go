package local

import (
	"encoding/json"
	"fmt"
	"github.com/meroxa/valve"
	"io/ioutil"
	"reflect"
	"time"
	"unsafe"
)

type Valve struct {
	fixturesPath string
}

func New(fixturesPath string) Valve {
	return Valve{fixturesPath: fixturesPath}
}

func (v Valve) Resources(name string) (valve.Resource, error) {
	return Resource{
		Name:         name,
		fixturesPath: mapFixturesPath(name, v.fixturesPath),
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
	Name string
	fixturesPath string
}

func (r Resource) Records(collection string, cfg valve.ResourceConfigs) (valve.Records, error) {
	return ReadFixtures(r.fixturesPath, collection)
}

func (r Resource) Write(rr valve.Records, collection string, cfg valve.ResourceConfigs) error {
	return nil
}

func ReadFixtures(path, collection string) (valve.Records, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return valve.Records{}, err
	}

	var records map[string]map[string]valve.Payload
	err = json.Unmarshal(b, &records)
	if err != nil {
		return valve.Records{}, err
	}

	var rr []valve.Record
	for k, r := range records[collection] {
		rr = append(rr, wrapRecord(k, r))
	}

	return valve.NewRecords(rr), nil
}

func mapFixturesPath(name, path string) string {
	return fmt.Sprintf("%s/%s.json", path, name)
}

func wrapRecord(key string, m map[string]interface{}) valve.Record {
	return valve.Record{
		Key: key,
		Payload: m,
		Timestamp: time.Now(),
	}
}
