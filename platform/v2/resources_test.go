package v2

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/meroxa/turbine-go"
)

func TestResources(t *testing.T) {

}

func TestListResources(t *testing.T) {
	tests := []struct {
		description string
		resources   []turbine.Resource
		output      []string
	}{
		{
			description: "Successfully list multiple resources",
			resources:   []turbine.Resource{testFunction{}, testFunction{}},
			output:      []string{"my-func", "my-func2"},
		},
		{
			description: "Successfully list zero resources",
			resources:   []turbine.Resource{},
			output:      nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			testTurbine := Turbine{
				resources: tc.resources,
			}
			list := testTurbine.ListFunctions()
			require.Equal(t, list, tc.output)
		})
	}
}

func TestRecords(t *testing.T) {

}

func TestWriteWithConfig(t *testing.T) {

}