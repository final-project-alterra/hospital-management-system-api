package data

import (
	"github.com/final-project-alterra/hospital-management-system-api/features/admins"
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	CreatedBy *uint
	UpdatedBy *uint

	Creating []Admin `gorm:"foreignkey:CreatedBy"`
	Updating []Admin `gorm:"foreignkey:UpdatedBy"`

	Email     string `gorm:"type:varchar(100);unique_index"`
	Password  string `gorm:"type:varchar(100)"`
	Name      string `gorm:"type:varchar(100)"`
	Phone     string `gorm:"type:varchar(100)"`
	Gender    string `gorm:"type:varchar(1)"`
	BirthDate string `gorm:"type:date"`
	Address   string
	ImageUrl  string
}

func (a Admin) ToAdminCore() admins.AdminCore {
	var createdBy int
	var updatedBy int

	if a.CreatedBy != nil {
		createdBy = int(*a.CreatedBy)
	}
	if a.UpdatedBy != nil {
		updatedBy = int(*a.UpdatedBy)
	}

	return admins.AdminCore{
		ID:        int(a.ID),
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
		CreatedBy: createdBy,
		UpdatedBy: updatedBy,
		Email:     a.Email,
		Password:  a.Password,
		Name:      a.Name,
		BirthDate: a.BirthDate,
		ImageUrl:  a.ImageUrl,
		Phone:     a.Phone,
		Address:   a.Address,
		Gender:    a.Gender,
	}
}

func ToSliceAdminCore(a []Admin) []admins.AdminCore {
	result := make([]admins.AdminCore, len(a))
	for i, v := range a {
		result[i] = v.ToAdminCore()
	}
	return result
}
