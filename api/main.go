package main

import (
	"config-generator/api/handlers"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.POST("/configure", handlers.ConfigurationHandler)
	e.POST("/service-mesh", handlers.ServiceMeshHandler)

	// old version
	//e.POST("/helm", handlers.HelmHandler)
	//e.POST("/cd-pipeline", handlers.CDPipelineHandler)

	e.Logger.Fatal(e.Start(":8081"))
}
