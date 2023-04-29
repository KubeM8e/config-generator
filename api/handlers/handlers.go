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
	user           = "Jon Doe"
	deploymentName = "Deployment-1"
	argoSuffix     = "-argo"
	helmSuffix     = "-helm"
)

func HelmHandler(c echo.Context) error {
	configsMap := make(map[string]interface{})

	// reads the JSON object and stores in the map
	err := json.NewDecoder(c.Request().Body).Decode(&configsMap)
	if err != nil {
		log.Fatal(err)
	}

	// save helm templates to a GitHub repo
	configs.DeploymentId = utils.ReadFromDB(user, deploymentName)
	utils.CreateGitHubRepo(configs.DeploymentId + helmSuffix)

	// generates values.yaml file from the map
	repoName := configs.DeploymentId + helmSuffix
	gitWorkTree, gitRepo := core.GenerateValuesYamlFile(configsMap, repoName)

	// generates helm charts
	core.ConfigureHelmChart(configsMap, gitWorkTree, gitRepo)

	return c.JSON(http.StatusOK, configsMap)
}

func CDPipelineHandler(c echo.Context) error {
	request := models.CDPipelineRequest{}

	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		log.Fatal(err)
	}

	// connect to database
	// TODO: move this to the correct place
	utils.ConnectMongoDB()
	configs.DeploymentId = utils.ReadFromDB(user, deploymentName)

	repoName := configs.DeploymentId + argoSuffix
	utils.CreateGitHubRepo(repoName)

	// configure application.yaml file
	response := configs.ConfigureCDPipeline(request, repoName)

	return c.JSON(http.StatusOK, &response)
}

func ServiceMeshHandler(c echo.Context) error {
	return nil
}
