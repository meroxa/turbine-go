package turbine

import (
	"context"
	"fmt"
	"os"

	"github.com/meroxa/turbine-go/pkg/proto/core"
)

func (tc *Turbine) RegisterSecret(name string) error {
	return tc.RegisterSecretWithContext(context.Background(), name)
}

// RegisterSecretWithContext pulls environment variables with the same name and ships them as Env Vars for functions
func (tc *Turbine) RegisterSecretWithContext(ctx context.Context, name string) error {
	val := os.Getenv(name)
	if val == "" {
		return fmt.Errorf("secret %q is invalid or not set", name)
	}

	_, err := tc.TurbineCore.RegisterSecret(ctx, &core.Secret{
		Name:  name,
		Value: val,
	})
	return err
}
