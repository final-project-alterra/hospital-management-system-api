package routes

import (
	"github.com/final-project-alterra/hospital-management-system-api/factory"
	"github.com/labstack/echo/v4"
)

func setupAuthRoutes(e *echo.Echo, presenter *factory.Presenter) {
	auth := e.Group("/auth")

	auth.POST("/login", presenter.AuthPresentation.PostLogin)
}
