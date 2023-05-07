package main

import (
	"config-generator/api/handlers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.CORS())

	e.POST("/configure/:id", handlers.ConfigurationHandler)
	e.POST("/service-mesh", handlers.ServiceMeshHandler)

	// old version
	//e.POST("/helm", handlers.HelmHandler)
	//e.POST("/cd-pipeline", handlers.CDPipelineHandler)

	e.Logger.Fatal(e.Start(":8081"))
}
