package routes

import (
	"github.com/final-project-alterra/hospital-management-system-api/factory"
	"github.com/labstack/echo/v4"
)

func setupNurseRoutes(e *echo.Echo, presenter *factory.Presenter) {
	// nurses := e.Group("/nurses", middleware.IsAdmin())
	nurses := e.Group("/nurses")

	nurses.GET("", presenter.NursePresentation.GetNurses)
	nurses.GET("/:nurseId", presenter.NursePresentation.GetDetailNurse)
	nurses.POST("", presenter.NursePresentation.PostNurse)
	nurses.PUT("", presenter.NursePresentation.PutEditNurse)
	nurses.PUT("/password", presenter.NursePresentation.PutEditNursePassword)
	nurses.DELETE("/:nurseId", presenter.NursePresentation.DeleteNurse)
}
