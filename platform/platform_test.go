package platform

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/meroxa-go/pkg/mock"
	"github.com/meroxa/turbine-go"
	"github.com/stretchr/testify/assert"
)

type TestFunc struct{}

func (t TestFunc) Process(r []turbine.Record) []turbine.Record {
	return []turbine.Record{}
}

func Test_Process(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name     string
		setupApp func() *Turbine
	}{
		{
			name: "without platform deployment",
			setupApp: func() *Turbine {
				return &Turbine{
					client: &Client{
						Client: mock.NewMockClient(ctrl),
					},
					functions: make(map[string]turbine.Function),
					resources: []turbine.Resource{},
					imageName: "image1",
					deploy:    false,
					config: turbine.AppConfig{
						Name:     "my-app",
						Pipeline: "my-pipe",
					},
					secrets: make(map[string]string),
					gitSha:  "sha123456789",
				}
			},
		},
		{
			name: "deploy on the platform",
			setupApp: func() *Turbine {
				c := mock.NewMockClient(ctrl)
				c.EXPECT().
					CreateFunction(
						gomock.Any(),
						&meroxa.CreateFunctionInput{
							Name:        "testfunc-sha12345",
							InputStream: "test",
							Image:       "image1",
							EnvVars:     make(map[string]string),
							Args:        []string{"testfunc"},
							Pipeline:    meroxa.PipelineIdentifier{Name: "my-pipe"},
						},
					).
					Return(&meroxa.Function{
						UUID: "1234-5678",
					}, nil).
					Times(1)

				return &Turbine{
					client:    &Client{Client: c},
					functions: make(map[string]turbine.Function),
					resources: []turbine.Resource{},
					imageName: "image1",
					deploy:    true,
					config: turbine.AppConfig{
						Name:     "my-app",
						Pipeline: "my-pipe",
					},
					secrets: make(map[string]string),
					gitSha:  "sha123456789",
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			app := tc.setupApp()
			app.Process(
				turbine.Records{
					Stream: "test",
				},
				TestFunc{},
			)
		})
	}
}

func Test_KafkaResourceWrite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name     string
		setupApp func() *Turbine
	}{
		{
			name: "write to Kafka resource",
			setupApp: func() *Turbine {
				c := mock.NewMockClient(ctrl)
				c.EXPECT().
					CreateConnector(
						gomock.Any(),
						&meroxa.CreateConnectorInput{
							ResourceName: "kafka1",
							PipelineName: "my-pipe",
							Configuration: map[string]interface{}{
								"conduit": "true",
								"topic":   "target-collection",
							},
							Type: "destination",
						}).
					Return(&meroxa.Connector{
						UUID: "1234-5678",
					}, nil).Times(1)

				c.EXPECT().GetPipelineByName(gomock.Any(), "my-pipe").
					Return(&meroxa.Pipeline{Name: "my-pipe", UUID: "1234-5678"}, nil).
					Times(1)

				c.EXPECT().GetResourceByNameOrID(gomock.Any(), "kafka1").
					Return(&meroxa.Resource{Name: "kafka1", UUID: "1234-5678", Type: "kafka"}, nil).
					Times(1)

				return &Turbine{
					client:    &Client{Client: c},
					functions: make(map[string]turbine.Function),
					resources: []turbine.Resource{},
					imageName: "image1",
					deploy:    true,
					config: turbine.AppConfig{
						Name:     "my-app",
						Pipeline: "my-pipe",
					},
					secrets: make(map[string]string),
					gitSha:  "sha123456789",
					appUUID: "1234-5678",
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			app := tc.setupApp()
			kafka1, _ := app.Resources("kafka1")
			kafka1.Write(turbine.Records{}, "target-collection")
		})
	}
}

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
