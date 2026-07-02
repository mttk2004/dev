package docker

import (
	"bytes"
	"os/exec"
	"strings"
)

// Container represents a local Docker container.
type Container struct {
	ID        string
	Name      string
	Status    string
	Image     string
	IsRunning bool
}

// IsDockerInstalled checks if docker command is available in PATH.
func IsDockerInstalled() bool {
	_, err := exec.LookPath("docker")
	return err == nil
}

// IsDockerRunning checks if the Docker daemon is running by calling 'docker info'.
func IsDockerRunning() bool {
	cmd := exec.Command("docker", "info")
	return cmd.Run() == nil
}

// GetContainers returns a list of all local docker containers.
func GetContainers() ([]Container, error) {
	cmd := exec.Command("docker", "ps", "-a", "--format", "{{.ID}}\t{{.Names}}\t{{.Status}}\t{{.Image}}")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	var containers []Container
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) < 4 {
			continue
		}

		id := parts[0]
		name := parts[1]
		status := parts[2]
		image := parts[3]
		isRunning := strings.HasPrefix(status, "Up")

		containers = append(containers, Container{
			ID:        id,
			Name:      name,
			Status:    status,
			Image:     image,
			IsRunning: isRunning,
		})
	}
	return containers, nil
}

// StartContainer starts a docker container.
func StartContainer(id string) error {
	return exec.Command("docker", "start", id).Run()
}

// StopContainer stops a running docker container.
func StopContainer(id string) error {
	return exec.Command("docker", "stop", id).Run()
}

// RestartContainer restarts a docker container.
func RestartContainer(id string) error {
	return exec.Command("docker", "restart", id).Run()
}

// GetContainerLogs returns the last 50 lines of logs for a docker container.
func GetContainerLogs(id string) (string, error) {
	cmd := exec.Command("docker", "logs", "--tail", "50", id)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}
