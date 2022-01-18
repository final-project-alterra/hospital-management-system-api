package data

import (
	"strings"

	"github.com/final-project-alterra/hospital-management-system-api/features/nurses"
	"gorm.io/gorm"
)

type Nurse struct {
	gorm.Model

	CreatedBy int
	UpdatedBy int

	Email     string `gorm:"type:varchar(64);unique;not null"`
	Password  string `gorm:"type:varchar(128);not null"`
	Name      string `gorm:"type:varchar(64);not null"`
	BirthDate string `gorm:"type:date;not null"`
	ImageUrl  string
	Phone     string `gorm:"type:varchar(16)"`
	Address   string
	Gender    string `gorm:"type:varchar(1);not null"`
}

func (n Nurse) ToNurseCore() nurses.NurseCore {
	return nurses.NurseCore{
		ID:        int(n.ID),
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,

		CreatedBy: n.CreatedBy,
		UpdatedBy: n.UpdatedBy,

		Email:     n.Email,
		Password:  n.Password,
		Name:      n.Name,
		BirthDate: strings.Split(n.BirthDate, "T")[0],
		ImageUrl:  n.ImageUrl,
		Phone:     n.Phone,
		Address:   n.Address,
		Gender:    n.Gender,
	}
}

func ToSliceNurseCore(n []Nurse) []nurses.NurseCore {
	result := make([]nurses.NurseCore, len(n))
	for i := range n {
		result[i] = n[i].ToNurseCore()
	}
	return result
}
