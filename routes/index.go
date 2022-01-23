package routes

import (
	"path"

	"github.com/final-project-alterra/hospital-management-system-api/factory"
	"github.com/final-project-alterra/hospital-management-system-api/middleware"
	"github.com/final-project-alterra/hospital-management-system-api/utils/project"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func SetupRoutes() *echo.Echo {
	e := echo.New()
	presenter := factory.New()

	e.Pre(echoMiddleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Static("/static", path.Join(project.GetMainDir(), "files"))

	setupAuthRoutes(e, presenter)

	setupAdminRoutes(e, presenter)

	setupNurseRoutes(e, presenter)

	setupPatientRoutes(e, presenter)

	setupDoctorRoutes(e, presenter)
	setupRoomRoutes(e, presenter)
	setupSpecialityRoutes(e, presenter)

	setupScheduleRoutes(e, presenter)

	setupOutpatientRoutes(e, presenter)

	return e
}
