package request

import "github.com/final-project-alterra/hospital-management-system-api/features/admins"

type CreateAdminRequest struct {
	CreatedBy int
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	Name      string `json:"name"`
	Age       int    `json:"age"`
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
		Age:       r.Age,
		Phone:     r.Phone,
		Address:   r.Address,
		Gender:    r.Gender,
	}
}
