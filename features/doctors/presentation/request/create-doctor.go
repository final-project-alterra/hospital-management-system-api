package request

import (
	"github.com/final-project-alterra/hospital-management-system-api/features/doctors"
)

type CreateDoctorRequest struct {
	CreatedBy    int
	SpecialityID int `json:"specialityId" validate:"gt=0"`
	RoomID       int `json:"roomId" validate:"gt=0"`

	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	Name      string `json:"name" validate:"required"`
	BirthDate string `json:"birthDate" validate:"required,ValidateBirthDate"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Gender    string `json:"gender" validate:"required,oneof='L' 'P'"`
}

func (d CreateDoctorRequest) ToDoctorCore() doctors.DoctorCore {
	return doctors.DoctorCore{
		CreatedBy: d.CreatedBy,

		Speciality: doctors.SpecialityCore{
			ID: d.SpecialityID,
		},
		Room: doctors.RoomCore{
			ID: d.RoomID,
		},

		Email:     d.Email,
		Password:  d.Password,
		Name:      d.Name,
		BirthDate: d.BirthDate,
		Phone:     d.Phone,
		Address:   d.Address,
		Gender:    d.Gender,
	}
}
