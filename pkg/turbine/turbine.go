package turbine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"os/exec"

	"github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"github.com/meroxa/turbine-core/pkg/client"
)

type Turbine interface {
	Resources(string) (Resource, error)
	ResourcesWithContext(context.Context, string) (Resource, error)

	Process(Records, Function) (Records, error)
	ProcessWithContext(context.Context, Records, Function) (Records, error)
}

type turbine struct {
	client.Client
}

func NewTurbineClient(ctx context.Context, turbineCoreAddress, gitSha, appPath string) (Turbine, error) {
	c, err := client.DialContext(ctx, turbineCoreAddress)
	if err != nil {
		return nil, err
	}

	tc := &turbine{Client: c}

	appName, err := appName(appPath)
	if err != nil {
		return nil, err
	}

	version, err := turbineGoVersion(ctx)
	if err != nil {
		return nil, err
	}

	req := core.InitRequest{
		AppName:        appName,
		ConfigFilePath: appPath,
		Language:       core.Language_GOLANG,
		GitSHA:         gitSha,
		TurbineVersion: version,
	}

	if _, err = tc.Init(ctx, &req); err != nil {
		return nil, err
	}

	return tc, nil
}

func appName(appPath string) (string, error) {
	b, err := os.ReadFile(appPath + "/" + "app.json")
	if err != nil {
		return "", err
	}

	type appConfig struct {
		Name string `json:"name"`
	}

	var ac appConfig
	if err = json.Unmarshal(b, &ac); err != nil {
		return "", err
	}

	if ac.Name == "" {
		return "", errors.New("application name is required to be specified in your app.json")
	}

	return ac.Name, nil
}

// turbineGoVersion will return the tag or hash of the turbine-go dependency of a given app.
func turbineGoVersion(ctx context.Context) (string, error) {
	var cmd *exec.Cmd

	cwd, _ := os.Getwd()
	cmd = exec.CommandContext(
		ctx,
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
