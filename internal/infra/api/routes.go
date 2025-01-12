package api

import (
	"github.com/labstack/echo/v4"
)

func (s *RestServer) SetUpRoutes() {
	s.UserRoutes()
}

func (s *RestServer) UserRoutes() {
	s.router.POST("/video", func(c echo.Context) error {
		return s.appHandler.videoProcessorHandler.Handle(c)
	})
}
