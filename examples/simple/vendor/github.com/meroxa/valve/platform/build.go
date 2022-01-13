package platform

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
	"io"
	"log"
	"os"

	"github.com/docker/docker/client"
)

func (v Valve) BuildDockerImage(name, path string) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("unable to init docker client; %s", err)
	}

	ctx := context.Background()
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
	//dockerFileTarReader := bytes.NewReader(buf.Bytes())

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

	tar, err := archive.TarWithOptions(".", &archive.TarOptions{
		Compression:     archive.Uncompressed,
		ExcludePatterns: []string{"simple", ".git", "fixtures"},
	})
	if err != nil {
		log.Fatalf("unable to create tar; %s", err)
	}

	resp, err := cli.ImageBuild(
		ctx,
		tar,
		types.ImageBuildOptions{
			Context:    tar,
			Dockerfile: "Dockerfile",
			Remove:     true,
			Tags:       []string{name}})
	if err != nil {
		log.Fatalf("unable to build docker image; %s", err)
	}
	defer resp.Body.Close()
	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		log.Fatalf("unable to to read image build response; %s", err)
	}
}
