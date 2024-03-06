package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

const (
	tomcatImageName = "tomcat:latest" // Adjust the Tomcat image name and tag as needed
)

func main() {
	// Connect to the local Docker daemon
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}

	// List running Docker containers
	containerList, err := listRunningContainers(cli)
	if err != nil {
		log.Fatal(err)
	}

	// Check if Tomcat is running in any container
	for _, container := range containerList {
		if isTomcatRunning(cli, container.ID) {
			log.Printf("Tomcat is running in container %s", container.ID[:10])
		}
	}
}

func listRunningContainers(cli *client.Client) ([]types.Container, error) {
	// Filter options to list only running containers
	options := types.ContainerListOptions{
		All:     false,
		Filters: map[string][]string{"status": {"running"}},
	}

	// Get the list of running containers
	containerList, err := cli.ContainerList(context.Background(), options)
	if err != nil {
		return nil, err
	}

	return containerList, nil
}

func isTomcatRunning(cli *client.Client, containerID string) bool {
	// Inspect the container to get information about its image
	containerInfo, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		log.Printf("Error inspecting container %s: %v", containerID[:10], err)
		return false
	}

	// Check if the container is running the Tomcat image
	if containerInfo.Config.Image == tomcatImageName {
		// Optionally, you can add more checks to verify the Tomcat process is running inside the container
		// For example, check logs for Tomcat startup messages
		logs, err := cli.ContainerLogs(context.Background(), containerID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
		if err != nil {
			log.Printf("Error getting container logs for %s: %v", containerID[:10], err)
			return false
		}
		defer logs.Close()

		// Read and print the container logs
		stdout, stderr := stdcopy.StdCopy(os.Stdout, os.Stderr, logs)
		log.Printf("Container logs for %s:\n%s\n%s", containerID[:10], stdout, stderr)

		// Return true if Tomcat is running, based on your specific criteria
		return true
	}

	return false
}
