package routes

import (
	"github.com/final-project-alterra/hospital-management-system-api/factory"
	"github.com/final-project-alterra/hospital-management-system-api/middleware"
	"github.com/labstack/echo/v4"
)

func setupSpecialityRoutes(e *echo.Echo, presenter *factory.Presenter) {
	speciality := e.Group("/specialities")

	speciality.GET("", presenter.DoctorPresentation.GetSpecialities, middleware.IsAuth())
	speciality.GET("/:specialityId", presenter.DoctorPresentation.GetDetailSpeciality, middleware.IsAuth())
	speciality.POST("", presenter.DoctorPresentation.PostSpeciality, middleware.IsAdmin())
	speciality.PUT("", presenter.DoctorPresentation.PutEditSpeciality, middleware.IsAdmin())
	speciality.DELETE("/:specialityId", presenter.DoctorPresentation.DeleteSpeciality, middleware.IsAdmin())
}
