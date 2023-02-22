package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func printImageList(ctx *context.Context, cli *client.Client) {

	images, err := cli.ImageList(*ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	for _, image := range images {
		fmt.Printf("%s %s\n", image.ID[:10], image.RepoTags)
	}
}

func printContainerList(ctx *context.Context, cli *client.Client) {

	containers, err := cli.ContainerList(*ctx, types.ContainerListOptions{All: true})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
}

func GetImage() {
	// Get from a remote repository
	// Get from docker host (local)
}

func BuildImage(ctx *context.Context, cli *client.Client, dockerFilePath string, imageName string) {
	cli.ImageBuild(*ctx, nil, types.ImageBuildOptions{})

}

func RunContainer(ctx *context.Context, cli *client.Client, imageName string) (string, error) {
	/*
		out, err := cli.ImagePull(*ctx, imageName, types.ImagePullOptions{})
		if err != nil {
			panic(err)
		}
		defer out.Close()
		io.Copy(os.Stdout, out)
	*/

	resp, err := cli.ContainerCreate(*ctx, &container.Config{
		Image: imageName}, nil, nil, nil, "lisa-container")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created container ID: %s. Starting it now...\n", resp.ID)

	if err := cli.ContainerStart(*ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Printf("Container started. ID: %s\n", resp.ID)

	statusC, errC := cli.ContainerWait(*ctx, resp.ID, "" /*container.WaitConditionNotRunning*/)
	select {
	case err := <-errC:
		if err != nil {
			log.Fatal(err)
		}
	case status := <-statusC:
		fmt.Printf("status.StatusCode: %#+v\n", status.StatusCode)
	}

	fmt.Println("Container finished.")

	return resp.ID, nil
}

func PrintContainerLogs(ctx *context.Context, cli *client.Client, containerID string) {
	fmt.Println("Printing container logs...")
	reader, err := cli.ContainerLogs(*ctx, containerID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Details: true})
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	bytes, err := io.Copy(os.Stdout, reader)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Copied %d bytes from container logs.\n", bytes)
}

func main() {
	// Context provides a means of transmitting deadlines, caller cancellations, and other request-scoped values across API boundaries and between processes.
	// Incoming requests to a server should create a Context, and outgoing calls to servers should accept a Context.
	// The chain of function calls between them must propagate the Context
	ctx := context.Background()

	// Initializes a new API client with a default HTTP client, and connects to the default Docker daemon socket
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	fmt.Println("Printing image list:")
	printImageList(&ctx, cli)

	fmt.Println("Printing container list:")
	printContainerList(&ctx, cli)

	containerId, err := RunContainer(&ctx, cli, "lisa-agent:latest")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Returned container ID: %s\n", containerId)

	PrintContainerLogs(&ctx, cli, containerId)

	cli.ContainerRemove(ctx, containerId, types.ContainerRemoveOptions{Force: true})

}
