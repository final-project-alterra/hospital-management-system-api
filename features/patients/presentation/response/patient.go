package response

import "github.com/final-project-alterra/hospital-management-system-api/features/patients"

type PatientResponse struct {
	ID      int    `json:"id"`
	NIK     string `json:"nik"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Gender  string `json:"gender"`
}

func DetailPatient(p patients.PatientCore) PatientResponse {
	return PatientResponse{
		ID:      p.ID,
		NIK:     p.NIK,
		Name:    p.Name,
		Age:     p.Age,
		Phone:   p.Phone,
		Address: p.Address,
		Gender:  p.Gender,
	}
}

func ListPatients(p []patients.PatientCore) []PatientResponse {
	result := make([]PatientResponse, len(p))
	for i := range p {
		result[i] = DetailPatient(p[i])
	}
	return result
}
