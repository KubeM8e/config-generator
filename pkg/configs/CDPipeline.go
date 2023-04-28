package configs

import (
	"config-generator/models"
	"config-generator/pkg/utils"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

// TODO: make them envs?
const (
	isSelfHeal    = true
	isPrune       = true
	helmDirectory = "helm"
	tmpDirectory  = "tmp"
)

var DeploymentId string

func ConfigureCDPipeline(configObject models.CDPipelineRequest) *models.ArgoCDApplicationConfig {

	applicationConfig := models.ArgoCDApplicationConfig{
		APIVersion: "argoproj.io/v1alpha1",
		Kind:       "Application",
	}

	applicationConfig.Metadata.Name = configObject.ApplicationName
	applicationConfig.Metadata.Namespace = "argocd" // the cluster must have a namespace called argocd
	applicationConfig.Spec.Project = "default"
	applicationConfig.Spec.Source.RepoURL = configObject.GitHubRepoURL
	applicationConfig.Spec.Source.TargetRevision = configObject.TargetRevision
	applicationConfig.Spec.Source.Path = configObject.Path
	applicationConfig.Spec.Destination.Server = configObject.ClusterURL
	applicationConfig.Spec.Destination.Namespace = configObject.ApplicationNamespace
	applicationConfig.Spec.SyncPolicy.SyncOptions = []string{"CreateNamespace=true"}
	applicationConfig.Spec.SyncPolicy.Automated.SelfHeal = isSelfHeal
	applicationConfig.Spec.SyncPolicy.Automated.Prune = isPrune

	gitWorkTree, gitRepo := utils.CloneGitHubRepo(DeploymentId)

	applicationYamlFile, err := os.Create(tmpDirectory + "/application.yaml")
	if err != nil {
		log.Fatal(err)
	}

	applicationYamlData, err := yaml.Marshal(&applicationConfig)
	if err != nil {
		log.Fatal(err)
	}

	_, err = applicationYamlFile.Write(applicationYamlData)
	if err != nil {
		log.Fatal(err)
	}

	utils.PushToGitHub(gitWorkTree, gitRepo)

	return &applicationConfig

}
