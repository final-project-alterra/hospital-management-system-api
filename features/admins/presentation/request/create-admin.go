package request

import (
	"github.com/final-project-alterra/hospital-management-system-api/features/admins"
)

type CreateAdminRequest struct {
	CreatedBy int
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	Name      string `json:"name" validate:"required"`
	BirthDate string `json:"birthDate" validate:"required,ValidateBirthDate"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Gender    string `json:"gender" validate:"required,oneof='L' 'P'"`
}

func (r CreateAdminRequest) ToAdminCore() admins.AdminCore {
	return admins.AdminCore{
		CreatedBy: r.CreatedBy,
		Email:     r.Email,
		Password:  r.Password,
		Name:      r.Name,
		BirthDate: r.BirthDate,
		Phone:     r.Phone,
		Address:   r.Address,
		Gender:    r.Gender,
	}
}
