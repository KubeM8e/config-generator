package configs

import (
	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	"os"
	"testing"
)

func TestCreatePrometheusStackValuesYaml(t *testing.T) {
	// set up a temporary git repository to test against
	repoPath, err := os.MkdirTemp("", "CreatePrometheusStackValuesYaml")
	require.NoError(t, err)
	defer os.RemoveAll(repoPath)
	repo, err := git.PlainInit(repoPath, false)
	require.NoError(t, err)
	_, err = repo.Worktree()
	require.NoError(t, err)

	// create mock data
	clusterIPs := []string{"192.168.0.1", "192.168.0.2"}
	repoName := "test-repo"

	// call the function being tested
	resultWorktree, resultRepo := CreatePrometheusStackValuesYaml(clusterIPs, repoName)

	// check that the function returns the expected results
	assert.NotNil(t, resultWorktree)
	assert.NotNil(t, resultRepo)

	// check that the file was created and pushed to the repo
	expectedFile := "monitoringValues.yaml"
	file, err := os.ReadFile(repoPath + "/" + expectedFile)
	require.NoError(t, err)
	assert.NotEmpty(t, file)

	// check that the values were updated as expected
	var values map[string]interface{}
	err = yaml.Unmarshal(file, &values)
	require.NoError(t, err)
	assert.Contains(t, values, "endpoints")
	endpoints := values["endpoints"].([]interface{})
	assert.Equal(t, clusterIPs[0], endpoints[0])
	assert.Equal(t, clusterIPs[1], endpoints[1])
}
