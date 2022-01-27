package platform

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/client"
)

func (v Valve) BuildDockerImage(name, path string) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("unable to init docker client; %s", err)
	}

	ctx := context.Background()

	// Generate Dockerfile internally
	//buf := new(bytes.Buffer)
	//tw := tar.NewWriter(buf)
	//defer tw.Close()
	//
	//dockerFileBytes := generateDockerfile()
	//
	//tarHeader := &tar.Header{
	//	Name: "Dockerfile",
	//	Size: int64(len(dockerFileBytes)),
	//}
	//err = tw.WriteHeader(tarHeader)
	//if err != nil {
	//	log.Fatalf("unable to write tar header; %s", err)
	//}
	//_, err = tw.Write(dockerFileBytes)
	//if err != nil {
	//	log.Fatalf("unable to write tar body; %s", err)
	//}
	//tar := bytes.NewReader(buf.Bytes())

	//resp, err := cli.ImageBuild(
	//	ctx,
	//	dockerFileTarReader,
	//	types.ImageBuildOptions{
	//		Context:    dockerFileTarReader,
	//		Dockerfile: "Dockerfile",
	//		Remove:     true,
	//		Tags:       []string{name}})
	//if err != nil {
	//	log.Fatalf("unable to build docker image; %s", err)
	//}
	//defer resp.Body.Close()
	//_, err = io.Copy(os.Stdout, resp.Body)
	//if err != nil {
	//	log.Fatalf("unable to to read image build response; %s", err)
	//}

	// Read local Dockerfile
	tar, err := archive.TarWithOptions(".", &archive.TarOptions{
		Compression:     archive.Uncompressed,
		ExcludePatterns: []string{"simple", ".git", "fixtures"},
	})
	if err != nil {
		log.Fatalf("unable to create tar; %s", err)
	}

	buildOptions := types.ImageBuildOptions{
		Context:    tar,
		Dockerfile: "Dockerfile",
		Remove:     true,
		Tags:       []string{imageName(name)}}

	resp, err := cli.ImageBuild(
		ctx,
		tar,
		buildOptions,
	)
	if err != nil {
		log.Fatalf("unable to build docker image; %s", err)
	}
	defer resp.Body.Close()
	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		log.Fatalf("unable to to read image build response; %s", err)
	}
}

func (v Valve) PushDockerImage(name string) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("unable to init docker client; %s", err)
	}
	authConfig := getAuthConfig()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	opts := types.ImagePushOptions{RegistryAuth: authConfig}
	rd, err := cli.ImagePush(ctx, imageName(name), opts)
	if err != nil {
		log.Fatalf("unable to push docker image; %s", err)
	}

	defer rd.Close()

	_, err = io.Copy(os.Stdout, rd)
	if err != nil {
		log.Fatalf("unable to to read image build response; %s", err)
	}
}

func getAuthConfig() string {
	dhUsername := os.Getenv("DOCKER_HUB_USERNAME")
	dhPassword := os.Getenv("DOCKER_HUB_PASSWORD")
	authConfig := types.AuthConfig{
		Username:      dhUsername,
		Password:      dhPassword,
		ServerAddress: "https://index.docker.io/v1/",
	}
	authConfigBytes, _ := json.Marshal(authConfig)
	return base64.URLEncoding.EncodeToString(authConfigBytes)
}

func imageName(name string) string {
	scope := os.Getenv("DOCKER_HUB_USERNAME")
	return strings.Join([]string{scope, name}, "/")
}

func generateDockerfile() []byte {
	return []byte(`
FROM golang:1.17 as build-env

WORKDIR /go/src/app
COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go install --tags=platform -mod=vendor ./examples/simple/...

FROM gcr.io/distroless/static
USER nonroot:nonroot
WORKDIR /app
COPY --from=build-env /go/bin/* /app
COPY --from=build-env /go/src/app/examples/simple/app.json /app
ENTRYPOINT ["/app/simple", "--serve"]
`)
}
