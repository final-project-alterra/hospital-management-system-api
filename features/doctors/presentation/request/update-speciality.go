package request

import "github.com/final-project-alterra/hospital-management-system-api/features/doctors"

type UpdateSpecialityRequest struct {
	ID   int    `json:"id" validate:"gt=0"`
	Name string `json:"name" validate:"required"`
}

func (s UpdateSpecialityRequest) ToSpecialityCore() doctors.SpecialityCore {
	return doctors.SpecialityCore{
		ID:   s.ID,
		Name: s.Name,
	}
}
