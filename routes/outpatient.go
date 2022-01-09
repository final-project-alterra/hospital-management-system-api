package routes

import (
	"github.com/final-project-alterra/hospital-management-system-api/factory"
	"github.com/final-project-alterra/hospital-management-system-api/middleware"
	"github.com/labstack/echo/v4"
)

func setupOutpatientRoutes(e *echo.Echo, presenter *factory.Presenter) {
	outpatients := e.Group("/outpatients")

	outpatients.GET("", presenter.SchedulePresentation.GetOutpatients, middleware.IsAuth())
	outpatients.GET("/:outpatientId", presenter.SchedulePresentation.GetDetailOutpatient, middleware.IsAuth())
	outpatients.POST("", presenter.SchedulePresentation.PostOutpatient, middleware.IsAdmin())
	outpatients.PUT("", presenter.SchedulePresentation.PutEditOutpatient, middleware.IsAdmin())
	outpatients.PUT("/cancel", presenter.SchedulePresentation.PutCancelOutpatient, middleware.IsAdmin())
	outpatients.PUT("/examine", presenter.SchedulePresentation.PutExamineOutpatient, middleware.IsAdmin())
	outpatients.PUT("/finish", presenter.SchedulePresentation.PutFinishOutpatient, middleware.IsAdmin())
	outpatients.DELETE("/:outpatientId", presenter.SchedulePresentation.DeleteOutpatient, middleware.IsAdmin())
}
