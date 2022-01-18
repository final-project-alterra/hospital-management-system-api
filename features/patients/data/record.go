package data

import (
	"strings"

	"github.com/final-project-alterra/hospital-management-system-api/features/patients"
	"gorm.io/gorm"
)

type Patient struct {
	gorm.Model
	CreatedBy int
	UpdatedBy int
	NIK       string `gorm:"unique_index;not null"`
	Name      string `gorm:"type:varchar(64);not null"`
	Phone     string `gorm:"type:varchar(16)"`
	Gender    string `gorm:"type:varchar(1);not null"`
	BirthDate string `gorm:"type:date;not null"`
	Address   string
}

func (p Patient) toPatientCore() patients.PatientCore {
	return patients.PatientCore{
		ID:        int(p.ID),
		CreatedBy: p.CreatedBy,
		UpdatedBy: p.UpdatedBy,
		NIK:       p.NIK,
		Name:      p.Name,
		BirthDate: strings.Split(p.BirthDate, "T")[0],
		Phone:     p.Phone,
		Address:   p.Address,
		Gender:    p.Gender,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func toSlicePatientCore(p []Patient) []patients.PatientCore {
	result := make([]patients.PatientCore, len(p))
	for i := range p {
		result[i] = p[i].toPatientCore()
	}
	return result
}
