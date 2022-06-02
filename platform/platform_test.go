package platform

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/meroxa/turbine-go"
)

func TestListResources(t *testing.T) {
	testCases := []struct {
		name                  string
		resourceCollection    []string
		expectedResourceNames []string
	}{
		{
			name:                  "ListResources returns resource names if several resources have been registered",
			resourceCollection:    []string{"piston", "nozzle"},
			expectedResourceNames: []string{"piston", "nozzle"},
		},
		{
			name:                  "ListResources returns resource names if a single resource has been registered",
			resourceCollection:    []string{"cylinder"},
			expectedResourceNames: []string{"cylinder"},
		},
		{
			name:                  "ListResources returns an empty list if no resources have been registered",
			resourceCollection:    []string{},
			expectedResourceNames: []string(nil),
		},
	}

	// Test setup
	// ==========
	// 1. mock Turbine client with a simplified app configuration handler
	origReadAppConfig := turbine.ReadAppConfig
	turbine.ReadAppConfig = func(appPath string) (turbine.AppConfig, error) {
		return turbine.AppConfig{}, nil
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 2. create a new Turbine client to make methods available for testing
			turbineMock := New(false, "engine")
			// 3. configure Turbine mock client with sample resources
			for _, name := range tc.resourceCollection {
				turbineMock.resources[name] = Resource{}
			}

			// Test execution
			// ==============
			result := turbineMock.ListResources()
			assert.ElementsMatch(t, tc.expectedResourceNames, result)

		})
	}

	// 4. reset Turbine client configuration handler
	turbine.ReadAppConfig = origReadAppConfig
}
