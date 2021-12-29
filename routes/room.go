package routes

import (
	"github.com/final-project-alterra/hospital-management-system-api/factory"
	"github.com/labstack/echo/v4"
)

func setupRoomRoutes(e *echo.Echo, presenter *factory.Presenter) {
	// room := e.Group("/rooms", middleware.Isroom())
	room := e.Group("/rooms")

	room.GET("", presenter.DoctorPresentation.GetRooms)
	room.POST("", presenter.DoctorPresentation.PostRoom)
	room.PUT("", presenter.DoctorPresentation.PutEditRoom)
	room.DELETE("/:roomId", presenter.DoctorPresentation.DeleteRoom)
}
