package main

import (
	"context"

	"github.com/Zyprush18/imagez/internal/routes"
	"github.com/labstack/echo/v5"
)

func main() {
	api := echo.New()
	vhost := map[string]*echo.Echo{}

	routes.Routes(api, vhost)

	sc := echo.StartConfig{
		Address: ":8000",
	}

	e:= echo.NewVirtualHostHandler(vhost) 

	if err:= sc.Start(context.Background(), e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
