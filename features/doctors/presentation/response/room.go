package response

import (
	"time"

	"github.com/final-project-alterra/hospital-management-system-api/features/doctors"
)

type RoomResponse struct {
	ID        int       `json:"id"`
	Floor     string    `json:"floor"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func DetailRoom(r doctors.RoomCore) RoomResponse {
	return RoomResponse{
		ID:        r.ID,
		Floor:     r.Floor,
		Code:      r.Code,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func ListRooms(r []doctors.RoomCore) []RoomResponse {
	rooms := make([]RoomResponse, len(r))
	for i := range r {
		rooms[i] = DetailRoom(r[i])
	}
	return rooms
}
