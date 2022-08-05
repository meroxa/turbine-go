package platform

import (
	"testing"

	"github.com/meroxa/turbine-go"
	"github.com/stretchr/testify/assert"
)

func TestListResources(t *testing.T) {
	testCases := []struct {
		name               string
		resourceCollection []string
		turbineResource    map[string]*Resource
	}{
		{
			name:               "ListResources returns resource names if several resources have been registered",
			resourceCollection: []string{"nozzle", "piston"},
			turbineResource: map[string]*Resource{
				"nozzle": {
					Name:        "nozzle",
					Collection:  "test",
					Source:      false,
					Destination: true,
				},
				"piston": {
					Name:        "piston",
					Collection:  "test 123",
					Source:      true,
					Destination: false,
				},
			},
		},
		{
			name:               "ListResources returns resource names if a single resource has been registered",
			resourceCollection: []string{"cylinder"},
			turbineResource: map[string]*Resource{
				"cylinder": {
					Name:        "cylinder",
					Collection:  "test",
					Source:      false,
					Destination: true,
				},
			},
		},
		{
			name:               "ListResources returns an empty list if no resources have been registered",
			resourceCollection: []string{},
			turbineResource:    map[string]*Resource{},
		},
	}

	// Test setup
	// ==========
	// 1. mock Turbine client with a simplified app configuration handler
	origReadAppConfig := turbine.ReadAppConfig
	turbine.ReadAppConfig = func(appName, appPath string) (turbine.AppConfig, error) {
		return turbine.AppConfig{
			Name: appName,
		}, nil
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 2. create a new Turbine client to make methods available for testing
			turbineMock := New(false, "engine", "app", "7c7f63ca-e906-4d0a-a488-65d8dbad1c89")
			// 3. configure Turbine mock client with sample resources
			for name := range tc.turbineResource {
				turbineMock.resources = append(turbineMock.resources, tc.turbineResource[name])
			}
			// Test execution
			// ==============
			result, err := turbineMock.ListResources()
			if err != nil {
				t.Errorf("no error expected; got %s", err.Error())
			}

			if len(result) != len(tc.turbineResource) {
				t.Errorf("incorrect number of resources returned")
			}

			assert.Equal(t, turbineMock.config.Name, "app")

			for _, el := range result {
				assert.Equal(t, tc.turbineResource[el.Name].Name, el.Name)
				assert.Equal(t, tc.turbineResource[el.Name].Collection, el.Collection)
				assert.Equal(t, tc.turbineResource[el.Name].Source, el.Source)
				assert.Equal(t, tc.turbineResource[el.Name].Destination, el.Destination)

			}
		})
	}

	// 4. reset Turbine client configuration handler
	turbine.ReadAppConfig = origReadAppConfig
}
