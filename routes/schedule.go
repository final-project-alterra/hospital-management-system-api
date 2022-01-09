package routes

import (
	"github.com/final-project-alterra/hospital-management-system-api/factory"
	"github.com/final-project-alterra/hospital-management-system-api/middleware"
	"github.com/labstack/echo/v4"
)

func setupScheduleRoutes(e *echo.Echo, presenter *factory.Presenter) {
	schedule := e.Group("/work-schedules")

	schedule.GET("", presenter.SchedulePresentation.GetWorkSchedules, middleware.IsAuth())
	schedule.POST("", presenter.SchedulePresentation.PostWorkSchedules, middleware.IsAdmin())
	schedule.PUT("", presenter.SchedulePresentation.PutEditWorkSchedule, middleware.IsAdmin())
	schedule.DELETE("/:workScheduleId", presenter.SchedulePresentation.DeleteWorkSchedule, middleware.IsAdmin())

	schedule.GET("/:workScheduleId/outpatients", presenter.SchedulePresentation.GetWorkScheduleOutpatients, middleware.IsAuth())
}
