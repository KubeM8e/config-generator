package core

import (
	"config-generator/models"
	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigureTemplatesAndCharts(t *testing.T) {
	// Set up input data for the test
	serviceMap := make(map[string]models.Microservice)
	microservice := models.Microservice{
		ServiceName: "myservice",
		AvgReplicas: 2,
		Envs:        nil,
		DockerImage: "docker.io/myuser/myservice:latest",
		MaxCPU:      "100m",
		MaxMemory:   "128Mi",
	}

	serviceMap["myservice"] = microservice

	appDataRequest := models.ConfigurationRequest{
		Microservices: serviceMap,
	}

	// Set up mock objects
	mockRepo := &git.Repository{}
	mockWorkTree := &git.Worktree{}

	// Set up expectations
	expectedWorktree := &git.Worktree{}
	expectedRepo := &git.Repository{}

	// Call the function being tested
	actualWorktree, actualRepo := ConfigureTemplatesAndCharts(appDataRequest, "my-repo-name", mockWorkTree, mockRepo)

	// Assert the results
	assert.Equal(t, expectedWorktree, actualWorktree)
	assert.Equal(t, expectedRepo, actualRepo)
}
