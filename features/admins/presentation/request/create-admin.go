package request

import (
	"time"

	"github.com/final-project-alterra/hospital-management-system-api/features/admins"
	"github.com/go-playground/validator/v10"
)

type CreateAdminRequest struct {
	CreatedBy int
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	Name      string `json:"name"`
	BirthDate string `json:"birthDate" validate:"required,ValidateCreateAdminBirthDate"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Gender    string `json:"gender" validate:"oneof='L' 'P'"`
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

func ValidateCreateAdminBirthDate(fl validator.FieldLevel) bool {
	input, ok := fl.Parent().Interface().(EditAdminRequest)
	if !ok {
		return false
	}

	_, err := time.Parse("2006-01-02", input.BirthDate)
	return err == nil
}
