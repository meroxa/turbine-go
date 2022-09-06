package v2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGoVersion(t *testing.T) {
	_, err := getGoVersion()
	require.NoError(t, err)
}
