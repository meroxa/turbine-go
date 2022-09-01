package v2

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/meroxa/turbine-go"
	"github.com/meroxa/turbine-go/platform"
)

type Turbine struct {
	client      *platform.Client
	functions   map[string]turbine.Function
	resources   []turbine.Resource
	deploy      bool
	deploySpec  *deploySpec
	specVersion string
	imageName   string
	appName     string
	config      turbine.AppConfig
	secrets     map[string]string
	gitSha      string
}

type deploySpec struct {
	Secrets    map[string]string `json:"secrets,omitempty"`
	Connectors []specConnector   `json:"connectors"`
	Functions  []specFunction    `json:"functions,omitempty"`
	Definition specDefinition    `json:"definition"`
}

type specConnector struct {
	Type       string                 `json:"type"`
	Resource   string                 `json:"resource"`
	Collection string                 `json:"collection"`
	Config     map[string]interface{} `json:"config,omitempty"`
}

type specFunction struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type specDefinition struct {
	AppName  string       `json:"app_name"`
	GitSha   string       `json:"git_sha"`
	Metadata specMetadata `json:"turbine"`
}

type specMetadata struct {
	Turbine     specTurbine `json:"turbine"`
	SpecVersion string      `json:"spec_version"`
}

type specTurbine struct {
	Language string `json:"language"`
	Version  string `json:"version"`
}

func New(deploy bool, imageName, appName, gitSha, spec string) *Turbine {
	c, err := platform.NewClient()
	if err != nil {
		log.Fatalln(err)
	}

	ac, err := turbine.ReadAppConfig(appName, "")
	if err != nil {
		log.Fatalln(err)
	}
	return &Turbine{
		client:      c,
		functions:   make(map[string]turbine.Function),
		resources:   []turbine.Resource{},
		imageName:   imageName,
		appName:     appName,
		deploy:      deploy,
		deploySpec:  &deploySpec{},
		specVersion: spec,
		config:      ac,
		secrets:     make(map[string]string),
		gitSha:      gitSha,
	}
}

func (t *Turbine) HandleSpec() (string, error) {
	t.deploySpec.Secrets = t.secrets

	version, err := getGoVersion()
	if err != nil {
		return "", err
	}

	t.deploySpec.Definition = specDefinition{
		AppName: t.appName,
		GitSha:  t.gitSha,
		Metadata: specMetadata{
			Turbine: specTurbine{
				Language: "golang",
				Version:  version,
			},
			SpecVersion: t.specVersion,
		},
	}

	bytes, err := json.MarshalIndent(t.deploySpec, "", "    ")
	// @TODO send deployment spec to Platform API if a deployment
	return string(bytes), err
}

func getGoVersion() (string, error) {
	cmd := exec.Command("go", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("unable to determine go version: %s", string(output))
	}
	words := strings.Split(string(output), " ")
	if len(words) < 3 {
		return "", fmt.Errorf("unable to determine go version: unexpected output %s", string(output))
	}
	version := words[2]
	version = strings.ReplaceAll(version, "go", "")
	return version, nil
}
