package response

import (
	"time"

	"github.com/final-project-alterra/hospital-management-system-api/features/patients"
)

type PatientResponse struct {
	ID        int       `json:"id"`
	NIK       string    `json:"nik"`
	Name      string    `json:"name"`
	BirthDate string    `json:"birthDate"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	Gender    string    `json:"gender"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func DetailPatient(p patients.PatientCore) PatientResponse {
	return PatientResponse{
		ID:        p.ID,
		NIK:       p.NIK,
		Name:      p.Name,
		BirthDate: p.BirthDate,
		Phone:     p.Phone,
		Address:   p.Address,
		Gender:    p.Gender,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func ListPatients(p []patients.PatientCore) []PatientResponse {
	result := make([]PatientResponse, len(p))
	for i := range p {
		result[i] = DetailPatient(p[i])
	}
	return result
}
