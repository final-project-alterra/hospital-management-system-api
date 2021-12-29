package response

import "github.com/final-project-alterra/hospital-management-system-api/features/doctors"

type RoomResponse struct {
	ID    int    `json:"id"`
	Floor string `json:"floor"`
	Code  string `json:"code"`
}

func DetailRoom(r doctors.RoomCore) RoomResponse {
	return RoomResponse{
		ID:    r.ID,
		Floor: r.Floor,
		Code:  r.Code,
	}
}

func ListRooms(r []doctors.RoomCore) []RoomResponse {
	rooms := make([]RoomResponse, len(r))
	for i := range r {
		rooms[i] = DetailRoom(r[i])
	}
	return rooms
}
