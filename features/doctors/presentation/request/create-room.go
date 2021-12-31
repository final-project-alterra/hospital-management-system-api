package request

import "github.com/final-project-alterra/hospital-management-system-api/features/doctors"

type CreateRoomRequest struct {
	Floor string `json:"floor" validate:"required"`
	Code  string `json:"code" validate:"required"`
}

func (r CreateRoomRequest) ToRoomCore() doctors.RoomCore {
	return doctors.RoomCore{
		Floor: r.Floor,
		Code:  r.Code,
	}
}
