package local

import (
	"encoding/json"
	"fmt"
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	_ "github.com/apache/beam/sdks/v2/go/pkg/beam/io/filesystem/local"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/io/textio"
	"github.com/meroxa/turbine-go"
	"github.com/meroxa/turbine-go/models"
	"os"
	"path"
	"time"
)

func (r Runner) Resource(name string) (turbine.Resource, error) {
	return Resource{
		scope:        r.scope,
		fixturesPath: r.config.Resources[name],
	}, nil
}

var _ turbine.Resource = (*Resource)(nil)

type Resource struct {
	scope        beam.Scope
	fixturesPath string
}

func (r Resource) Records(collection string) (models.Records, error) {
	binPath, err := os.Executable()
	if err != nil {
		return models.Records{}, err
	}
	if r.fixturesPath == "" {
		return models.Records{},
			fmt.Errorf("must specify fixtures path to data for source resources in order to run locally")
	}
	dirPath := path.Dir(binPath)
	pwd := fmt.Sprintf("%s/%s", dirPath, r.fixturesPath)
	rr, err := readFixtures(pwd, collection)
	if err != nil {
		return models.Records{}, err
	}

	var stringRecs []string
	for _, v := range models.GetRecords(rr) {
		stringRecs = append(stringRecs, v.String())
	}

	pcol := beam.CreateList(r.scope, stringRecs)
	if err != nil {
		return models.Records{}, err
	}
	return models.NewRecordsWithPcol(pcol), nil
}

func (r Resource) RecordsWithConfig(collection string, cfg models.ResourceConfig) (models.Records, error) {
	return r.Records(collection)
}

func (r Resource) Write(records models.Records, collection string) error {
	textio.Write(r.scope, "out.txt", models.GetPcol(records))
	return nil
}

func (r Resource) WriteWithConfig(records models.Records, collection string, cfg models.ResourceConfig) error {
	return r.Write(records, collection)
}

type fixtureRecord struct {
	Key       string
	Value     map[string]interface{}
	Timestamp string
}

func readFixtures(path, collection string) (models.Records, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return models.Records{}, err
	}

	var records map[string][]fixtureRecord
	err = json.Unmarshal(b, &records)
	if err != nil {
		return models.Records{}, err
	}

	var rr []models.Record
	for _, r := range records[collection] {
		rr = append(rr, wrapRecord(r))
	}

	return models.NewRecords(rr), nil
}

func wrapRecord(m fixtureRecord) models.Record {
	b, _ := json.Marshal(m.Value)

	var t time.Time
	if m.Timestamp == "" {
		t = time.Now()
	} else {
		t, _ = time.Parse(time.RFC3339, m.Timestamp)
	}

	return models.Record{
		Key:       m.Key,
		Payload:   b,
		Timestamp: t,
	}
}
