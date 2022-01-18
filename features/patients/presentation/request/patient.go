package request

import "github.com/final-project-alterra/hospital-management-system-api/features/patients"

type CreatePatientRequest struct {
	CreatedBy int
	NIK       string `json:"nik" validate:"required"`
	Name      string `json:"name" validate:"required"`
	Phone     string `json:"phone"`
	BirthDate string `json:"birthDate" validate:"required,ValidateBirthDate"`
	Address   string `json:"address"`
	Gender    string `json:"gender" validate:"required,oneof='L' 'P"`
}

type UpdatePatientRequest struct {
	UpdatedBy int
	ID        int    `json:"id" validate:"gt=0"`
	Name      string `json:"name" validate:"required"`
	Phone     string `json:"phone"`
	BirthDate string `json:"age" validate:"required,ValidateBirthDate"`
	Address   string `json:"address"`
	Gender    string `json:"gender" validate:"required,oneof='L' 'P"`
}

func (p CreatePatientRequest) ToPatientCore() patients.PatientCore {
	return patients.PatientCore{
		CreatedBy: p.CreatedBy,
		NIK:       p.NIK,
		Name:      p.Name,
		BirthDate: p.BirthDate,
		Phone:     p.Phone,
		Address:   p.Address,
		Gender:    p.Gender,
	}
}

func (p UpdatePatientRequest) ToPatientCore() patients.PatientCore {
	return patients.PatientCore{
		ID:        p.ID,
		UpdatedBy: p.UpdatedBy,
		Name:      p.Name,
		BirthDate: p.BirthDate,
		Phone:     p.Phone,
		Address:   p.Address,
		Gender:    p.Gender,
	}
}
