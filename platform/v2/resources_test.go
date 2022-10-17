package v2

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/meroxa/turbine-core/pkg/ir"
	"github.com/meroxa/turbine-go"
	"github.com/meroxa/turbine-go/platform"
)

func TestResources(t *testing.T) {
	testTurbine := Turbine{}

	for count := 0; count < 3; count++ {
		require.Equal(t, count, len(testTurbine.resources))

		_, err := testTurbine.Resources(fmt.Sprintf("res%d", count))
		require.NoError(t, err)
	}
}

func TestListResources(t *testing.T) {
	tests := []struct {
		description string
		resources   []turbine.Resource
		output      []platform.ResourceWithCollection
	}{
		{
			description: "Successfully list multiple resources",
			resources: []turbine.Resource{
				&Resource{Name: "my-res", Source: true, Destination: false, Collection: "source-table"},
				&Resource{Name: "my-res2", Source: false, Destination: true, Collection: "destination-table"}},
			output: []platform.ResourceWithCollection{
				{
					Source: true, Destination: false, Name: "my-res", Collection: "source-table",
				},
				{
					Source: false, Destination: true, Name: "my-res2", Collection: "destination-table",
				},
			},
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
			list, err := testTurbine.ListResources()
			require.NoError(t, err)
			require.Equal(t, list, tc.output)
		})
	}
}

func TestRecords(t *testing.T) {
	type parameters struct {
		Collection string
		Cfg        turbine.ConnectionOptions
	}

	tests := []struct {
		description string
		input       parameters
		err         error
	}{
		{
			description: "Fail to add a source connector with an empty collection",
			input: parameters{
				Collection: "",
				Cfg:        turbine.ConnectionOptions{{Field: "a", Value: "b"}, {Field: "c", Value: "d"}},
			},
			err: fmt.Errorf("please provide a collection name to Records()"),
		},
		{
			description: "Successfully add one source connector",
			input: parameters{
				Collection: "collection1",
				Cfg:        turbine.ConnectionOptions{{Field: "a", Value: "b"}, {Field: "c", Value: "d"}},
			},
			err: nil,
		},
		{
			description: "Fail to add a second source connector",
			input: parameters{
				Collection: "collection2",
				Cfg:        turbine.ConnectionOptions{{Field: "a", Value: "b"}, {Field: "c", Value: "d"}},
			},
			err: fmt.Errorf("only one call to Records() is allowed per Meroxa Data Application"),
		},
	}

	testResource := Resource{v: &Turbine{deploySpec: &ir.DeploymentSpec{}}}
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			numConnectors := len(testResource.v.deploySpec.Connectors)
			_, err := testResource.Records(tc.input.Collection, tc.input.Cfg)
			if err != nil {
				require.NotNil(t, tc.err)
				require.Equal(t, err, tc.err)
			} else {
				require.NoError(t, tc.err)

				l := len(testResource.v.deploySpec.Connectors)
				require.Equal(t, l, numConnectors+1)
				latestConnector := testResource.v.deploySpec.Connectors[l-1]
				require.Equal(t, latestConnector.Collection, tc.input.Collection)
				require.Equal(t, latestConnector.Config, tc.input.Cfg.ToMap())
			}
		})
	}
}

func TestWrite(t *testing.T) {
	tests := []struct {
		description string
		collections []string
		err         error
	}{
		{
			description: "Fail to add a destination connector with an empty collection",
			collections: []string{""},
			err:         fmt.Errorf("please provide a collection name to Write()"),
		},
		{
			description: "Successfully add multiple destination connectors",
			collections: []string{"collection1", "collection2"},
			err:         nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			testResource := Resource{v: &Turbine{deploySpec: &ir.DeploymentSpec{}}}
			for i, collection := range tc.collections {
				err := testResource.Write(turbine.Records{}, collection)
				if err != nil {
					require.NotNil(t, tc.err)
					require.Equal(t, err, tc.err)
				} else {
					require.NoError(t, tc.err)

					l := len(testResource.v.deploySpec.Connectors)
					require.Equal(t, l, i+1)
					latestConnector := testResource.v.deploySpec.Connectors[l-1]
					require.Equal(t, latestConnector.Collection, collection)
					require.Equal(t, latestConnector.Config, map[string]interface{}{})
				}
			}
		})
	}

}

func TestWriteWithConfig(t *testing.T) {
	type parameters struct {
		Collection string
		Cfg        turbine.ConnectionOptions
	}

	tests := []struct {
		description string
		input       []parameters
		err         error
	}{
		{
			description: "Fail to add a destination connector with an empty collection",
			input:       []parameters{{"", nil}},
			err:         fmt.Errorf("please provide a collection name to WriteWithConfig()"),
		},
		{
			description: "Successfully add multiple destination connectors",
			input: []parameters{
				{"collection1",
					turbine.ConnectionOptions{{Field: "a", Value: "b"}, {Field: "c", Value: "d"}},
				},
				{"collection2",
					turbine.ConnectionOptions{{Field: "e", Value: "f"}, {Field: "g", Value: "h"}},
				},
			},
			err: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			testResource := Resource{v: &Turbine{deploySpec: &ir.DeploymentSpec{}}}
			for i, input := range tc.input {
				err := testResource.WriteWithConfig(turbine.Records{}, input.Collection, input.Cfg)
				if err != nil {
					require.NotNil(t, tc.err)
					require.Equal(t, err, tc.err)
				} else {
					require.NoError(t, tc.err)

					l := len(testResource.v.deploySpec.Connectors)
					require.Equal(t, l, i+1)
					latestConnector := testResource.v.deploySpec.Connectors[l-1]
					require.Equal(t, latestConnector.Collection, input.Collection)
					require.Equal(t, latestConnector.Config, input.Cfg.ToMap())
				}
			}
		})
	}

}
