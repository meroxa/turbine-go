package v2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
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

	version, err := getTurbineGoVersion()
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

// getTurbineGoVersion will return the tag or hash of the turbine-go dependency of a given app.
func getTurbineGoVersion() (string, error) {
	var cmd *exec.Cmd

	cwd, _ := os.Getwd()
	cmd = exec.CommandContext(
		context.TODO(),
		"go",
		"list", "-m", "-f", "'{{ .Version }}'", "github.com/meroxa/turbine-go")
	fmtErr := fmt.Errorf(
		"unable to determine the version of turbine-go used by the Meroxa Application at %s",
		cwd)

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmtErr
	}

	version := strings.TrimSpace(string(stdout))
	chars := []rune(version)
	if chars[0] == 'v' {
		// Looks like v0.0.0-20221024132549-e6470e58b719
		const sections = 3
		parts := strings.Split(version, "-")
		if len(parts) < sections {
			return "", fmtErr
		}
		version = parts[2]
	}
	return version, nil
}
