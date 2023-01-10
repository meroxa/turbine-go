package app

import (
	"context"
	"fmt"
	"os"

	"github.com/meroxa/turbine-go/turbine/core"
)

// RegisterSecret pulls environment variables with the same name and ships them as Env Vars for functions
func (t *Turbine) RegisterSecret(name string) error {
	val := os.Getenv(name)
	if val == "" {
		return fmt.Errorf("secret %q is invalid or not set", name)
	}

	_, err := t.TurbineServiceClient.RegisterSecret(context.TODO(), &core.Secret{
		Name:  name,
		Value: val,
	})
	return err
}
