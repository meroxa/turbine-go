package v2

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/meroxa/turbine-core/pkg/ir"
	"github.com/meroxa/turbine-go"
	"github.com/meroxa/turbine-go/platform"
)

type Turbine struct {
	client      *platform.Client
	functions   map[string]turbine.Function
	resources   []turbine.Resource
	deploy      bool
	deploySpec  *ir.DeploymentSpec
	specVersion string
	imageName   string
	appName     string
	config      turbine.AppConfig
	secrets     map[string]string
	gitSha      string
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
		deploySpec:  &ir.DeploymentSpec{},
		specVersion: spec,
		config:      ac,
		secrets:     make(map[string]string),
		gitSha:      gitSha,
	}
}

func (t *Turbine) DeploymentSpec() (string, error) {
	t.deploySpec.Secrets = t.secrets

	version, err := getGoVersion()
	if err != nil {
		return "", err
	}

	t.deploySpec.Definition = ir.DefinitionSpec{
		GitSha: t.gitSha,
		Metadata: ir.MetadataSpec{
			Turbine: ir.TurbineSpec{
				Language: ir.GoLang,
				Version:  version,
			},
			SpecVersion: t.specVersion,
		},
	}

	bytes, err := json.Marshal(t.deploySpec)
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
