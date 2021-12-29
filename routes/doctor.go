package routes

import (
	"github.com/final-project-alterra/hospital-management-system-api/factory"
	"github.com/labstack/echo/v4"
)

func setupDoctorRoutes(e *echo.Echo, presenter *factory.Presenter) {
	// doctor := e.Group("/doctors", middleware.IsAdmin())
	doctor := e.Group("/doctors")

	doctor.GET("", presenter.DoctorPresentation.GetDoctors)
	doctor.GET("/:doctorId", presenter.DoctorPresentation.GetDetailDoctor)
	doctor.POST("", presenter.DoctorPresentation.PostDoctor)
	doctor.PUT("", presenter.DoctorPresentation.PutEditDoctor)
	doctor.PUT("/password", presenter.DoctorPresentation.PutEditDoctorPassword)
	doctor.DELETE("/:doctorId", presenter.DoctorPresentation.DeleteDoctor)
}
