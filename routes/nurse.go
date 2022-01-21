package routes

import (
	"github.com/final-project-alterra/hospital-management-system-api/factory"
	"github.com/final-project-alterra/hospital-management-system-api/middleware"
	"github.com/labstack/echo/v4"
)

func setupNurseRoutes(e *echo.Echo, presenter *factory.Presenter) {
	nurses := e.Group("/nurses")

	nurses.GET("", presenter.NursePresentation.GetNurses, middleware.IsAuth())
	nurses.GET("/:nurseId", presenter.NursePresentation.GetDetailNurse, middleware.IsAuth())
	nurses.POST("", presenter.NursePresentation.PostNurse, middleware.IsAdmin())
	nurses.PUT("", presenter.NursePresentation.PutEditNurse, middleware.IsAdmin())
	nurses.PUT("/password", presenter.NursePresentation.PutEditNursePassword, middleware.IsAdmin())
	nurses.PUT("/image-profile", presenter.NursePresentation.PutEditImageProfile, middleware.IsAdmin())
	nurses.DELETE("/:nurseId", presenter.NursePresentation.DeleteNurse, middleware.IsAdmin())
	nurses.DELETE("/:nurseId/image-profile", presenter.NursePresentation.DeleteImageProfile, middleware.IsAdmin())

	nurses.GET("/:nurseId/work-schedules", presenter.SchedulePresentation.GetNurseWorkSchedules, middleware.IsAuth())
}
