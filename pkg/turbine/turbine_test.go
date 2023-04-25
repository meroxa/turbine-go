package turbine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTurbineGoVersion(t *testing.T) {
	ctx := context.Background()
	_, err := turbineGoVersion(ctx)
	require.NoError(t, err)
}
