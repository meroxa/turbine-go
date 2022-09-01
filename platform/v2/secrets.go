package v2

import (
	"errors"
	"os"
)

// RegisterSecret pulls environment variables with the same name and ships them as Env Vars for functions
func (t *Turbine) RegisterSecret(name string) error {
	val := os.Getenv(name)
	if val == "" {
		return errors.New("secret is invalid or not set")
	}
	t.secrets[name] = val
	return nil
}
