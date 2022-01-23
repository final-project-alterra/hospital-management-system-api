package request

import "github.com/final-project-alterra/hospital-management-system-api/features/nurses"

type CreateNurseRequest struct {
	CreatedBy int

	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	Name      string `json:"name" validate:"required"`
	BirthDate string `json:"birthDate" validate:"required,ValidateBirthDate"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Gender    string `json:"gender" validate:"required,oneof='L' 'P'"`
}

func (c CreateNurseRequest) ToCore() nurses.NurseCore {
	return nurses.NurseCore{
		CreatedBy: c.CreatedBy,
		Email:     c.Email,
		Password:  c.Password,
		Name:      c.Name,
		BirthDate: c.BirthDate,
		Phone:     c.Phone,
		Address:   c.Address,
		Gender:    c.Gender,
	}
}
