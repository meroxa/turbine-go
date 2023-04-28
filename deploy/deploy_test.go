package deploy_test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/meroxa/turbine-core/pkg/ir"

	"github.com/meroxa/turbine-go/deploy"
	"github.com/stretchr/testify/assert"
)

func Test_CreateDockerfile_OldSpec(t *testing.T) {
	testCases := []struct {
		desc            string
		appName         string
		expectedAppName string
		pwd             string
	}{
		{
			desc:            "create dockerfile with provided app name",
			appName:         "myapp",
			expectedAppName: "myapp",
			pwd:             t.TempDir(),
		},
		{
			desc:            "create dockerfile with app name from json.file",
			expectedAppName: "expectedApp",
			pwd:             setupAppJson(t),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			dockerfile := path.Join(tc.pwd, "Dockerfile")
			err := deploy.CreateDockerfile(tc.appName, tc.pwd, ir.SpecVersion_0_1_1)
			if assert.NoError(t, err) && assert.FileExists(t, dockerfile) {
				v, err := os.ReadFile(dockerfile)
				if err != nil {
					t.Fatal(err)
				}
				assert.Contains(t, string(v), fmt.Sprintf("COPY %s.cross %s", tc.expectedAppName, tc.expectedAppName))
				assert.Contains(t, string(v), fmt.Sprintf("ENTRYPOINT [%q, %q]", "/app/"+tc.expectedAppName, "--serve"))
			}
		})
	}
}

func Test_CreateDockerfile(t *testing.T) {
	testCases := []struct {
		desc            string
		appName         string
		expectedAppName string
		pwd             string
	}{
		{
			desc:            "create dockerfile with provided app name",
			appName:         "myapp",
			expectedAppName: "myapp",
			pwd:             t.TempDir(),
		},
		{
			desc:            "create dockerfile with app name from json.file",
			expectedAppName: "expectedApp",
			pwd:             setupAppJson(t),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			dockerfile := path.Join(tc.pwd, "Dockerfile")
			err := deploy.CreateDockerfile(tc.appName, tc.pwd, ir.SpecVersion_0_2_0)
			if assert.NoError(t, err) && assert.FileExists(t, dockerfile) {
				v, err := os.ReadFile(dockerfile)
				if err != nil {
					t.Fatal(err)
				}

				assert.Contains(t, string(v), "FROM golang:1.20 as builder")
				assert.Contains(t, string(v), "WORKDIR /builder")
				assert.Contains(t, string(v), "COPY . .")
				assert.Contains(t, string(v), "ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64")
				assert.Contains(t, string(v), fmt.Sprintf("RUN go build -tags server -o %s ./...", tc.expectedAppName))

				assert.Contains(t, string(v), "FROM gcr.io/distroless/static")
				assert.Contains(t, string(v), "USER nobody")
				assert.Contains(t, string(v), "WORKDIR /app")
				assert.Contains(t, string(v), "COPY --from=builder /builder/app.json /app")
				assert.Contains(t, string(v), fmt.Sprintf("COPY --from=builder /builder/%s /app", tc.expectedAppName))

				entryPoint := fmt.Sprintf("ENTRYPOINT [%q, %q, %q]", path.Join("/app", tc.expectedAppName), "server", "-serve")
				assert.Contains(t, string(v), entryPoint)
			}
		})
	}
}

func setupAppJson(t *testing.T) string {
	tmpdir := t.TempDir()
	if err := os.WriteFile(
		path.Join(tmpdir, "app.json"),
		[]byte(`{
				  "name": "expectedApp",
				  "language": "golang",
				  "environment": "common",
				  "resources": {
				    "source_name": "fixtures/demo-cdc.json"
				  }
				}`),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	return tmpdir
}
