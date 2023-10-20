package build

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	sdk "github.com/meroxa/turbine-go/v2/pkg/turbine"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"github.com/meroxa/turbine-core/pkg/client"
)

var _ sdk.Turbine = (*builder)(nil)

type builder struct {
	client.Client
	runProcess bool
}

func NewBuildClient(ctx context.Context, turbineCoreAddress, gitSha, appPath string, runProcess bool) (*builder, error) {
	c, err := client.DialContext(ctx, turbineCoreAddress)
	if err != nil {
		return nil, err
	}

	b := &builder{
		Client:     c,
		runProcess: runProcess,
	}

	appName, err := appName(appPath)
	if err != nil {
		return nil, err
	}

	version, err := turbineGoVersion(ctx)
	if err != nil {
		return nil, err
	}

	req := pb.InitRequest{
		AppName:        appName,
		ConfigFilePath: appPath,
		Language:       pb.Language_GOLANG,
		GitSHA:         gitSha,
		TurbineVersion: version,
	}

	if _, err = b.Init(ctx, &req); err != nil {
		return nil, err
	}

	return b, nil
}

func appName(appPath string) (string, error) {
	b, err := os.ReadFile(filepath.Join(appPath, "app.json"))
	if err != nil {
		return "", err
	}

	ac := struct {
		Name string `json:"name"`
	}{}
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
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return "", fmt.Errorf("unable to determine turbine-go version")
	}

	parse := func(s string) string {
		v := strings.Split(s, "-")
		if len(v) < 3 {
			return s
		}
		return v[2]
	}

	for _, m := range bi.Deps {
		if m.Path == "github.com/meroxa/turbine-go/v2" { // this path is the same, regardless of OS
			return parse(m.Version), nil
		}
	}
	return "", fmt.Errorf("unable to find turbine-go in modules")
}
