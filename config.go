package turbine

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type AppConfig struct {
	Name        string            `json:"name"`
	Environment string            `json:"environment"`
	Pipeline    string            `json:"pipeline"`
	Resources   map[string]string `json:"resources"`
}

func ReadAppConfig() (AppConfig, error) {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("unable to locate executable path; error: %s", err)
	}

	projPath := path.Dir(exePath)
	b, err := ioutil.ReadFile(projPath + "/" + "app.json")
	if err != nil {
		return AppConfig{}, err
	}

	var ac AppConfig
	err = json.Unmarshal(b, &ac)
	if err != nil {
		return AppConfig{}, err
	}
	return ac, nil
}
