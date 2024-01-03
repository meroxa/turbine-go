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

	sdk "github.com/meroxa/turbine-go/v3/pkg/turbine"

	client "github.com/meroxa/turbine-core/v2/pkg/client"
	"github.com/meroxa/turbine-core/v2/proto/turbine/v2"
)

var _ sdk.Turbine = (*builder)(nil)

type builder struct {
	c          client.Client
	runProcess bool
}

func NewBuildClient(ctx context.Context, turbineCoreAddress, gitSha, appPath string, runProcess bool) (*builder, error) {
	c, err := client.DialContext(ctx, turbineCoreAddress)
	if err != nil {
		return nil, err
	}

	b := &builder{
		c:          c,
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

	req := turbinev2.InitRequest{
		AppName:        appName,
		ConfigFilePath: appPath,
		Language:       turbinev2.Language_GOLANG,
		GitSHA:         gitSha,
		TurbineVersion: version,
	}

	if _, err = b.c.Init(ctx, &req); err != nil {
		return nil, err
	}

	return b, nil
}

func (b *builder) Source(name, pluginName string, opts ...sdk.Option) (sdk.Source, error) {
	return b.SourceWithContext(context.Background(), name, pluginName, opts...)
}

func (b *builder) SourceWithContext(
	ctx context.Context,
	name string,
	pluginName string,
	opts ...sdk.Option,
) (sdk.Source, error) {
	var config sdk.OptionConfig

	config.Apply(opts...)

	resp, err := b.c.AddSource(ctx, &turbinev2.AddSourceRequest{
		Name: name,
		Plugin: &turbinev2.Plugin{
			Name:   pluginName,
			Config: config.PluginConfig,
		},
	})
	if err != nil {
		return nil, err
	}

	return &source{
		streamName: resp.StreamName,
		c:          b.c,
	}, nil
}

func (b *builder) Destination(name, pluginName string, opts ...sdk.Option) (sdk.Destination, error) {
	return b.DestinationWithContext(context.Background(), name, pluginName, opts...)
}

func (b *builder) DestinationWithContext(
	ctx context.Context,
	name string,
	pluginName string,
	opts ...sdk.Option,
) (sdk.Destination, error) {
	var config sdk.OptionConfig

	config.Apply(opts...)

	resp, err := b.c.AddDestination(ctx, &turbinev2.AddDestinationRequest{
		Name: name,
		Plugin: &turbinev2.Plugin{
			Name:   pluginName,
			Config: config.PluginConfig,
		},
	})
	if err != nil {
		return nil, err // todo wrap errors
	}

	return &destination{
		id: resp.Id,
		c:  b.c,
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

	parse := func(s string) string {
		v := strings.Split(s, "-")
		if len(v) < 3 {
			return s
		}
		return v[2]
	}

	for _, m := range bi.Deps {
		if m.Path == "github.com/meroxa/turbine-go/v3" { // this path is the same, regardless of OS
			return parse(m.Version), nil
		}
	}
	return "", fmt.Errorf("unable to find turbine-go in modules")
}
