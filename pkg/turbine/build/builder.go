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
}

func NewBuildClient(ctx context.Context, turbineCoreAddress, gitSha, appPath string) (*builder, error) {
	c, err := client.DialContext(ctx, turbineCoreAddress)
	if err != nil {
		return nil, err
	}

	b := &builder{Client: c}

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

func (b *builder) Resources(name string) (sdk.Resource, error) {
	return b.ResourcesWithContext(context.Background(), name)
}

func (b *builder) ResourcesWithContext(ctx context.Context, name string) (sdk.Resource, error) {
	r, err := b.GetResource(
		ctx,
		&pb.GetResourceRequest{
			Name: name,
		})
	if err != nil {
		return nil, err
	}

	return &resource{
		Resource: r,
		Client:   b.Client,
	}, nil
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

	parse := func(s string) (string, error) {
		v := strings.Split(s, "-")
		if len(v) < 3 {
			return "", fmt.Errorf("unable to parse version: %s", s)
		}
		return v[2], nil
	}

	for _, m := range bi.Deps {
		if m.Path == filepath.Join("github.com", "meroxa", "turbine-go") {
			return parse(m.Version)
		}
	}
	return "", fmt.Errorf("unable to find turbine-go in modules")
}
