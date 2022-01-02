package routes

import (
	"github.com/final-project-alterra/hospital-management-system-api/factory"
	"github.com/final-project-alterra/hospital-management-system-api/middleware"
	"github.com/labstack/echo/v4"
)

func setupRoomRoutes(e *echo.Echo, presenter *factory.Presenter) {
	room := e.Group("/rooms")

	room.GET("", presenter.DoctorPresentation.GetRooms, middleware.IsAuth())
	room.POST("", presenter.DoctorPresentation.PostRoom, middleware.IsAdmin())
	room.PUT("", presenter.DoctorPresentation.PutEditRoom, middleware.IsAdmin())
	room.DELETE("/:roomId", presenter.DoctorPresentation.DeleteRoom, middleware.IsAdmin())
}
