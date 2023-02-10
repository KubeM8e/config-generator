package main

import (
	"config-generator/api/handlers"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.POST("/core/charts", handlers.CoreChartsHandler)
	e.POST("/cd-pipeline", handlers.CDPipelineHandler)
	e.POST("/service-mesh", handlers.ServiceMeshHandler)
	e.Logger.Fatal(e.Start(":8081"))
}
