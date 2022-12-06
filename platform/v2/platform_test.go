package v2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGoVersion(t *testing.T) {
	_, err := getTurbineGoVersion()
	require.NoError(t, err)
}
