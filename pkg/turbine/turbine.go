//go:generate mockgen -source=turbine.go -package=mock -destination=mock/turbine_mock.go TurbineCore

package turbine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/meroxa/turbine-go/pkg/proto/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TurbineCore interface {
	core.TurbineServiceClient
}

type Turbine struct {
	TurbineCore
}

func NewCoreServer(ctx context.Context, turbineCoreAddress, gitSha string) (*Turbine, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial(turbineCoreAddress, opts...)
	if err != nil {
		return nil, err
	}

	tc := &Turbine{
		TurbineCore: core.NewTurbineServiceClient(conn),
	}

	if err := tc.Initialize(ctx, gitSha); err != nil {
		return nil, err
	}

	return tc, nil
}

func (tc *Turbine) Initialize(ctx context.Context, gitSha string) error {
	path, err := appPath()
	if err != nil {
		return err
	}

	appName, err := appName(path)
	if err != nil {
		return err
	}

	version, err := turbineGoVersion(ctx)
	if err != nil {
		return err
	}

	req := core.InitRequest{
		AppName:        appName,
		ConfigFilePath: path,
		Language:       core.Language_GOLANG,
		GitSHA:         gitSha,
		TurbineVersion: version,
	}

	_, err = tc.Init(ctx, &req)
	return err
}

func appPath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("unable to locate executable path; error: %s", err)
	}
	return path.Dir(exePath), nil
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
