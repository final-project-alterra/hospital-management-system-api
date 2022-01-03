package routes

import (
	"github.com/final-project-alterra/hospital-management-system-api/factory"
	"github.com/final-project-alterra/hospital-management-system-api/middleware"
	"github.com/labstack/echo/v4"
)

func setupPatientRoutes(e *echo.Echo, presenter *factory.Presenter) {
	patient := e.Group("/patients")

	patient.GET("", presenter.PatientPresentation.GetPatients, middleware.IsAuth())
	patient.GET("/:patientId", presenter.PatientPresentation.GetDetailPatient, middleware.IsAuth())
	patient.POST("", presenter.PatientPresentation.PostPatient, middleware.IsAdmin())
	patient.PUT("", presenter.PatientPresentation.PutEditPatient, middleware.IsAdmin())
	patient.DELETE("/:patientId", presenter.PatientPresentation.DeletePatient, middleware.IsAdmin())
}
