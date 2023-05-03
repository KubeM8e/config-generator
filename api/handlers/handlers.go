package handlers

import (
	"config-generator/models"
	"config-generator/pkg/configs"
	"config-generator/pkg/configs/core"
	"config-generator/pkg/utils"
	"encoding/json"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

const (
	argoSuffix = "-argo"
	helmSuffix = "-helm"
)

func ConfigurationHandler(c echo.Context) error {
	projectId := c.Param("id")
	request := models.ConfigurationRequest{}

	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		log.Printf("Could not decode the request: %s", err)
	}

	// repos
	repoNameHelm := projectId + helmSuffix
	repoNameArgo := projectId + argoSuffix
	helmRepoURL := "https://github.com/Shenali-SJ/" + repoNameHelm + ".git" // monitoringValues.yaml is also added to same repo
	argoRepoURL := "https://github.com/Shenali-SJ/" + repoNameArgo + ".git"

	// creates database
	//utils.CreateAppDataDB(request.AppName, request.Version)
	// reads from database
	//appId := utils.ReadFromDB(request.AppName)

	// creates helm GitHub repo
	utils.CreateGitHubRepo(repoNameHelm)

	// generates values.yaml
	gitWorkTree, gitRepo := core.GenerateValueAndChartFiles(request, repoNameHelm)

	// creates argo GitHub repo
	utils.CreateGitHubRepo(repoNameArgo)

	// generates helm templates
	argoWorktree, argoRepo := core.ConfigureTemplatesAndCharts(request, repoNameArgo, gitWorkTree, gitRepo)

	configs.ConfigureCDPipeline(request.AppName, helmRepoURL, request.ClusterURL, repoNameArgo, argoWorktree, argoRepo)

	response := models.RepoResponse{
		HelmRepo: helmRepoURL,
		ArgoRepo: argoRepoURL,
	}

	return c.JSON(http.StatusOK, &response)
}

//func HelmHandler(c echo.Context) error {
//	configsMap := make(map[string]interface{})
//
//	// reads the JSON object and stores in the map
//	err := json.NewDecoder(c.Request().Body).Decode(&configsMap)
//	if err != nil {
//		log.Printf("Could not decode the request: %s", err)
//	}
//
//	// save helm templates to a GitHub repo
//	appId := utils.ReadFromDB2(user, deploymentName)
//	utils.CreateGitHubRepo(appId + helmSuffix)
//
//	// generates values.yaml file from the map
//	repoName := appId + helmSuffix
//	gitWorkTree, gitRepo := core.GenerateValuesYamlFile2(configsMap, repoName)
//
//	// generates helm charts
//	core.ConfigureHelmChart2(configsMap, gitWorkTree, gitRepo)
//
//	return c.JSON(http.StatusOK, configsMap)
//}

//func CDPipelineHandler(c echo.Context) error {
//	request := models.CDPipelineRequest{}
//
//	err := json.NewDecoder(c.Request().Body).Decode(&request)
//	if err != nil {
//		log.Printf("Could not decode the request: %s", err)
//	}
//
//	// connect to database
//	// TODO: move this to the correct place
//	//configs.DeploymentId = utils.ReadFromDB(user, deploymentName)
//	// reads from database
//	//appId := utils.ReadFromDB(appName)
//
//	repoName := appId + argoSuffix
//	utils.CreateGitHubRepo(repoName)
//
//	// configure application.yaml file
//	response := configs.ConfigureCDPipeline(request, repoName)
//
//	return c.JSON(http.StatusOK, &response)
//}

func ServiceMeshHandler(c echo.Context) error {
	return nil
}
