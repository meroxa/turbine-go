package build

import (
	"context"
	"fmt"
	"os"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
)

func (b *builder) RegisterSecret(name string) error {
	return b.RegisterSecretWithContext(context.Background(), name)
}

// RegisterSecretWithContext pulls environment variables with the same name and ships them as Env Vars for functions
func (b *builder) RegisterSecretWithContext(ctx context.Context, name string) error {
	val := os.Getenv(name)
	if val == "" {
		return fmt.Errorf("secret %q is invalid or not set", name)
	}

	_, err := b.Client.RegisterSecret(ctx, &pb.Secret{
		Name:  name,
		Value: val,
	})
	return err
}
