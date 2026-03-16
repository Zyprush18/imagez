package routes

import (
	"github.com/Zyprush18/imagez/internal/handler"
	"github.com/Zyprush18/imagez/internal/service"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)



func Routes(e *echo.Echo, vhost map[string]*echo.Echo)  {
	e.Use(middleware.RequestLogger())
	
	vhost["api.localhost:8000"] = e
	apiV1 := e.Group("/v1")
	
	serviceImage := service.NewServiceImage()
	HandleImage :=  handler.NewHandleImage(serviceImage)

	apiV1.POST("/convert", HandleImage.Convert)
	apiV1.POST("/resize", HandleImage.Resize)
	apiV1.POST("/compress", HandleImage.Compress)
}