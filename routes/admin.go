package routes

import (
	"github.com/final-project-alterra/hospital-management-system-api/factory"
	"github.com/labstack/echo/v4"
)

func setupAdminRoutes(e *echo.Echo, presenter *factory.Presenter) {
	// admin := e.Group("/admins", middleware.IsAdmin())
	admin := e.Group("/admins")

	admin.GET("", presenter.AdminPresentation.GetAdmins)
	admin.GET("/:adminId", presenter.AdminPresentation.GetDetailAdmin)
	admin.POST("", presenter.AdminPresentation.PostCreateAdmin)
	admin.PUT("", presenter.AdminPresentation.PutEditAdmin)
	admin.PUT("/password", presenter.AdminPresentation.PutEditAdminPassword)
	admin.DELETE("/:adminId", presenter.AdminPresentation.DeleteAdmin)
}
