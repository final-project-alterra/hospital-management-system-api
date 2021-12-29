package request

import "github.com/final-project-alterra/hospital-management-system-api/features/doctors"

type UpdateRoomRequest struct {
	ID    int    `json:"id" validate:"required"`
	Floor string `json:"floor" validate:"required"`
	Code  string `json:"code" validate:"required"`
}

func (r UpdateRoomRequest) ToRoomCore() doctors.RoomCore {
	return doctors.RoomCore{
		ID:    r.ID,
		Floor: r.Floor,
		Code:  r.Code,
	}
}
