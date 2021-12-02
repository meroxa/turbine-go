package valve

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func ReadFixtures(resource, collection string) ([]Record, error) {
	path := mapFixturesPath(resource)
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err

	}
	defer jsonFile.Close()

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var records map[string]map[string]Payload
	err = json.Unmarshal(b, &records)

	var rr []Record
	for k, r := range records[collection] {
		rr = append(rr, wrapRecord(k, r))
	}

	return rr, err
}

func mapFixturesPath(name string) string {
	return fmt.Sprintf("./fixtures/%s.json", name)
}