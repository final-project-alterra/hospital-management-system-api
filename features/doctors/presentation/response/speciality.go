package response

import (
	"time"

	"github.com/final-project-alterra/hospital-management-system-api/features/doctors"
)

type SpecialityResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func DetailSpeciality(s doctors.SpecialityCore) SpecialityResponse {
	return SpecialityResponse{
		ID:        s.ID,
		Name:      s.Name,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func ListSpecialities(sc []doctors.SpecialityCore) []SpecialityResponse {
	result := make([]SpecialityResponse, len(sc))
	for i := range sc {
		result[i] = DetailSpeciality(sc[i])
	}
	return result
}
