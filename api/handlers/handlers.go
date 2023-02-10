package handlers

import (
	"config-generator/models"
	"config-generator/pkg/configs"
	"config-generator/pkg/utils"
	"encoding/json"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

const (
	user           = "Jon Doe"
	deploymentName = "Deployment-1"
)

func CoreChartsHandler(c echo.Context) error {
	return nil
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
	utils.CreateGitHubRepo(configs.DeploymentId)

	// configure application.yaml file
	response := configs.ConfigureCDPipeline(request)

	return c.JSON(http.StatusOK, &response)
}

func ServiceMeshHandler(c echo.Context) error {
	return nil
}
