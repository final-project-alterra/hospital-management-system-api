package response

import "github.com/final-project-alterra/hospital-management-system-api/features/doctors"

type SpecialityResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func DetailSpeciality(s doctors.SpecialityCore) SpecialityResponse {
	return SpecialityResponse{
		ID:   s.ID,
		Name: s.Name,
	}
}

func ListSpecialities(sc []doctors.SpecialityCore) []SpecialityResponse {
	result := make([]SpecialityResponse, len(sc))
	for i := range sc {
		result[i] = DetailSpeciality(sc[i])
	}
	return result
}
