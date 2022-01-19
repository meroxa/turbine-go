package valve

import (
	"encoding/json"
	"io/ioutil"
)

type AppConfig struct {
	Name        string            `json:"name"`
	Environment string            `json:"environment"`
	Pipeline    string            `json:"pipeline"`
	Resources   map[string]string `json:"resources"`
}

func ReadAppConfig() (AppConfig, error) {
	b, err := ioutil.ReadFile("app.json")
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
