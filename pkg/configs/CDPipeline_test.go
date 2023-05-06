package configs

import (
	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigureCDPipeline(t *testing.T) {
	appName := "myapp"
	repoURL := "https://github.com/myorg/myrepo.git"
	clusterURL := "https://mycluster.com"
	repoName := "myrepo"
	worktree := &git.Worktree{}     // create a dummy worktree
	repository := &git.Repository{} // create a dummy repository

	// Call the function and get the application config
	appConfig := ConfigureCDPipeline(appName, repoURL, clusterURL, repoName, worktree, repository)

	// Check that the application config was set up correctly
	assert.Equal(t, "argoproj.io/v1alpha1", appConfig.APIVersion)
	assert.Equal(t, "Application", appConfig.Kind)
	assert.Equal(t, appName, appConfig.Metadata.Name)
	assert.Equal(t, "argocd", appConfig.Metadata.Namespace)
	assert.Equal(t, "default", appConfig.Spec.Project)
	assert.Equal(t, repoURL, appConfig.Spec.Source.RepoURL)
	assert.Equal(t, "main", appConfig.Spec.Source.TargetRevision)
	assert.Equal(t, "./", appConfig.Spec.Source.Path)
	assert.Equal(t, []string{"values.yaml"}, appConfig.Spec.Source.Helm.ValueFiles)
	assert.Equal(t, clusterURL, appConfig.Spec.Destination.Server)
	assert.Equal(t, "argoapp", appConfig.Spec.Destination.Namespace)
	assert.Equal(t, []string{"CreateNamespace=true"}, appConfig.Spec.SyncPolicy.SyncOptions)
	assert.Equal(t, isSelfHeal, appConfig.Spec.SyncPolicy.Automated.SelfHeal)
	assert.Equal(t, isPrune, appConfig.Spec.SyncPolicy.Automated.Prune)

	// TODO: add more tests to check that the function handles errors correctly
}
