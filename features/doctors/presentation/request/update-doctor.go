package request

import "github.com/final-project-alterra/hospital-management-system-api/features/doctors"

type UpdateDoctorRequest struct {
	ID        int `json:"id" validate:"gt=0"`
	UpdatedBy int

	SpecialityID int `json:"specialityId" validate:"gt=0"`
	RoomID       int `json:"roomId" validate:"gt=0"`

	Name      string `json:"name" validate:"required"`
	BirthDate string `json:"birthDate" validate:"required,ValidateBirthDate"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Gender    string `json:"gender" validate:"required,oneof='L' 'P'"`
}

func (d UpdateDoctorRequest) ToDoctorCore() doctors.DoctorCore {
	return doctors.DoctorCore{
		ID:        d.ID,
		UpdatedBy: d.UpdatedBy,

		Speciality: doctors.SpecialityCore{
			ID: d.SpecialityID,
		},
		Room: doctors.RoomCore{
			ID: d.RoomID,
		},

		Name:      d.Name,
		BirthDate: d.BirthDate,
		Phone:     d.Phone,
		Address:   d.Address,
		Gender:    d.Gender,
	}
}
