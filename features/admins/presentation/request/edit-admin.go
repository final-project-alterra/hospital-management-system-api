package request

import (
	"github.com/final-project-alterra/hospital-management-system-api/features/admins"
)

type EditAdminRequest struct {
	ID        int `json:"id" validate:"required"`
	UpdatedBy int
	Name      string `json:"name" validate:"required"`
	BirthDate string `json:"birthDate" validate:"required,ValidateBirthDate"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Gender    string `json:"gender" validate:"required,oneof='L' 'P'"`
}

func (r EditAdminRequest) ToAdminCore() admins.AdminCore {
	return admins.AdminCore{
		ID:        r.ID,
		UpdatedBy: r.UpdatedBy,
		Name:      r.Name,
		BirthDate: r.BirthDate,
		Phone:     r.Phone,
		Address:   r.Address,
		Gender:    r.Gender,
	}
}
