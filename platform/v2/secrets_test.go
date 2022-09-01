package v2

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegisterSecret(t *testing.T) {
	tests := []struct {
		description string
		key         string
		value       string
		err         error
	}{
		{
			description: "Successfully register secret",
			key:         "SUPER_SECRET",
			value:       "deadbeef",
			err:         nil,
		},
		{
			description: "Fail to register empty secret",
			key:         "SUPER_SECRET",
			value:       "",
			err:         errors.New("secret is invalid or not set"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			testTurbine := Turbine{
				secrets: map[string]string{},
			}
			os.Setenv(tc.key, tc.value)

			err := testTurbine.RegisterSecret(tc.key)
			require.Equal(t, err, tc.err)
		})
	}
}
