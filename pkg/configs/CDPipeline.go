package configs

import (
	"config-generator/models"
	"config-generator/pkg/utils"
	"github.com/go-git/go-git/v5"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

// TODO: make them envs?
const (
	isSelfHeal      = true
	isPrune         = true
	applicationYaml = "application.yaml"
	TmpArgoFolder   = "argo"
)

//var DeploymentId string

func ConfigureCDPipeline(appName string, repoURL string, clusterURL string, repoName string, worktree *git.Worktree, repository *git.Repository) *models.ArgoCDApplicationConfig {

	var gitWorkTree *git.Worktree
	var gitRepo *git.Repository
	if worktree == nil {
		gitWorkTree, gitRepo = utils.CloneGitHubRepo(repoName, TmpArgoFolder)
	}

	applicationConfig := models.ArgoCDApplicationConfig{
		APIVersion: "argoproj.io/v1alpha1",
		Kind:       "Application",
	}

	applicationConfig.Metadata.Name = appName
	applicationConfig.Metadata.Namespace = "argocd" // the cluster must have a namespace called argocd
	applicationConfig.Spec.Project = "default"
	applicationConfig.Spec.Source.RepoURL = repoURL
	//applicationConfig.Spec.Source.TargetRevision = "HEAD"
	applicationConfig.Spec.Source.TargetRevision = "main"
	applicationConfig.Spec.Source.Path = "./" // base directory of GitHub repo
	applicationConfig.Spec.Source.Helm.ValueFiles = []string{"values.yaml"}
	applicationConfig.Spec.Destination.Server = clusterURL
	applicationConfig.Spec.Destination.Namespace = "argoapp"
	applicationConfig.Spec.SyncPolicy.SyncOptions = []string{"CreateNamespace=true"}
	applicationConfig.Spec.SyncPolicy.Automated.SelfHeal = isSelfHeal
	applicationConfig.Spec.SyncPolicy.Automated.Prune = isPrune

	//gitWorkTree, gitRepo := utils.CloneGitHubRepo(repoName, TmpArgoFolder)

	applicationYamlFile, err := os.Create(TmpArgoFolder + "/application.yaml")
	if err != nil {
		log.Printf("Could not create argo dir: %s", err)
	}

	applicationYamlData, err := yaml.Marshal(&applicationConfig)
	if err != nil {
		log.Printf("Could not marshal argo : %s", err)
	}

	_, err = applicationYamlFile.Write(applicationYamlData)
	if err != nil {
		log.Printf("Could not write argo : %s", err)
	}

	if worktree == nil {
		utils.PushToGitHub(gitWorkTree, gitRepo, []string{applicationYaml})
	} else {
		utils.PushToGitHub(worktree, repository, []string{applicationYaml})
	}

	return &applicationConfig

}

func ConfigureCDPipeline2(configObject models.CDPipelineRequest, repoName string) *models.ArgoCDApplicationConfig {

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

	gitWorkTree, gitRepo := utils.CloneGitHubRepo(repoName, TmpArgoFolder)

	applicationYamlFile, err := os.Create(TmpArgoFolder + "/application.yaml")
	if err != nil {
		log.Printf("Could not create argo dir: %s", err)
	}

	applicationYamlData, err := yaml.Marshal(&applicationConfig)
	if err != nil {
		log.Printf("Could not marshal argo : %s", err)
	}

	_, err = applicationYamlFile.Write(applicationYamlData)
	if err != nil {
		log.Printf("Could not write argo : %s", err)
	}

	utils.PushToGitHub(gitWorkTree, gitRepo, []string{applicationYaml})

	return &applicationConfig

}
