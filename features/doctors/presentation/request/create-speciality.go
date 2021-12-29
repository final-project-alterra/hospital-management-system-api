package request

import "github.com/final-project-alterra/hospital-management-system-api/features/doctors"

type CreateSpecialityRequest struct {
	Name string `json:"name" validate:"required"`
}

func (s CreateSpecialityRequest) ToSpecialityCore() doctors.SpecialityCore {
	return doctors.SpecialityCore{
		Name: s.Name,
	}
}
