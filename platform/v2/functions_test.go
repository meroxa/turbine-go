package v2

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/meroxa/turbine-core/pkg/ir"
	"github.com/meroxa/turbine-go"
)

type testFunction struct {
}

func (testFunction) Process(r []turbine.Record) []turbine.Record {
	return r
}

func TestGetFunction(t *testing.T) {
	tests := []struct {
		description string
		name        string
		functions   map[string]turbine.Function
		ok          bool
	}{
		{
			description: "Successfully find function by name",
			name:        "my-func",
			functions:   map[string]turbine.Function{"my-func": testFunction{}},
			ok:          true,
		},
		{
			description: "Fail to find non-existent function by name",
			name:        "my-func2",
			functions:   map[string]turbine.Function{},
			ok:          false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			testTurbine := Turbine{
				functions: tc.functions,
			}
			fun, ok := testTurbine.GetFunction(tc.name)
			require.Equal(t, ok, tc.ok)
			require.Equal(t, fun, tc.functions[tc.name])
		})
	}
}

func TestListFunctions(t *testing.T) {
	tests := []struct {
		description string
		functions   map[string]turbine.Function
		output      []string
	}{
		{
			description: "Successfully list multiple functions",
			functions:   map[string]turbine.Function{"my-func": testFunction{}, "my-func2": testFunction{}},
			output:      []string{"my-func", "my-func2"},
		},
		{
			description: "Successfully list zero functions",
			functions:   map[string]turbine.Function{},
			output:      nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			testTurbine := Turbine{
				functions: tc.functions,
			}
			list := testTurbine.ListFunctions()
			require.Equal(t, len(list), len(tc.output))
			for _, f := range tc.output {
				found := false
				for _, g := range list {
					if f == g {
						found = true
						break
					}
				}
				if !found {
					t.Fatalf("missing %s in ListFunctions outout", f)
				}
			}
		})
	}
}

func TestProcess(t *testing.T) {
	tests := []struct {
		description string
		image       string
		fun         turbine.Function
		input       turbine.Records
	}{
		{
			description: "Successfully Process testFunction",
			image:       "foo:bar",
			fun:         testFunction{},
			input:       turbine.Records{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			testTurbine := Turbine{
				imageName:  tc.image,
				functions:  make(map[string]turbine.Function),
				deploySpec: &ir.DeploymentSpec{},
			}
			output := testTurbine.Process(tc.input, tc.fun)
			require.Equal(t, tc.input, output)
			list := testTurbine.ListFunctions()
			require.Equal(t, len(list), 1)
			require.Equal(t, list[0], "testfunction")
		})
	}
}
