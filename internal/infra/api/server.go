package api

import (
	"hackaton-video-processor-worker/internal/infra/api/handlers"

	"github.com/labstack/echo/v4"
)

type RestServer struct {
	router     *echo.Echo
	appHandler *AppHandlers
}

type AppHandlers struct {
	videoProcessorHandler *handlers.VideoHandler
	// Add new handlers here
}

func NewRestService(router *echo.Echo, appHandler *AppHandlers) *RestServer {
	return &RestServer{
		router:     router,
		appHandler: appHandler,
	}
}

func SetUpServer() {
	router := echo.New()
	appHandler := configHandlers()

	server := NewRestService(router, appHandler)
	server.SetUpRoutes()

	server.router.Logger.Fatal(router.Start(":8080"))
}
