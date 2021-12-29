package routes

import (
	"github.com/final-project-alterra/hospital-management-system-api/factory"
	"github.com/labstack/echo/v4"
)

func setupSpecialityRoutes(e *echo.Echo, presenter *factory.Presenter) {
	// speciality := e.Group("/specialities", middleware.IsAdmin())
	speciality := e.Group("/specialities")

	speciality.GET("", presenter.DoctorPresentation.GetSpecialities)
	speciality.GET("/:specialityId", presenter.DoctorPresentation.GetDetailSpeciality)
	speciality.POST("", presenter.DoctorPresentation.PostSpeciality)
	speciality.PUT("", presenter.DoctorPresentation.PutEditSpeciality)
	speciality.DELETE("/:specialityId", presenter.DoctorPresentation.DeleteSpeciality)
}
