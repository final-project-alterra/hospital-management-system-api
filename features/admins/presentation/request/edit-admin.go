package request

import (
	"time"

	"github.com/final-project-alterra/hospital-management-system-api/features/admins"
	"github.com/go-playground/validator/v10"
)

type EditAdminRequest struct {
	ID        int `json:"id" validate:"required"`
	UpdatedBy int
	Name      string `json:"name"`
	BirthDate string `json:"birthDate" validate:"required,ValidateEditAdminBirthDate"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Gender    string `json:"gender"`
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

func ValidateEditAdminBirthDate(fl validator.FieldLevel) bool {
	input, ok := fl.Parent().Interface().(EditAdminRequest)
	if !ok {
		return false
	}

	_, err := time.Parse("2006-01-02", input.BirthDate)
	return err == nil
}
