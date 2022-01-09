package routes

import (
	"github.com/final-project-alterra/hospital-management-system-api/factory"
	"github.com/final-project-alterra/hospital-management-system-api/middleware"
	"github.com/labstack/echo/v4"
)

func setupDoctorRoutes(e *echo.Echo, presenter *factory.Presenter) {
	doctor := e.Group("/doctors")

	doctor.GET("", presenter.DoctorPresentation.GetDoctors, middleware.IsAuth())
	doctor.GET("/:doctorId", presenter.DoctorPresentation.GetDetailDoctor, middleware.IsAuth())
	doctor.POST("", presenter.DoctorPresentation.PostDoctor, middleware.IsAdmin())
	doctor.PUT("", presenter.DoctorPresentation.PutEditDoctor, middleware.IsAdmin())
	doctor.PUT("/password", presenter.DoctorPresentation.PutEditDoctorPassword, middleware.IsAdmin())
	doctor.DELETE("/:doctorId", presenter.DoctorPresentation.DeleteDoctor, middleware.IsAdmin())

	doctor.GET("/:doctorId/work-schedules", presenter.SchedulePresentation.GetDoctorWorkSchedules, middleware.IsAuth())
}
