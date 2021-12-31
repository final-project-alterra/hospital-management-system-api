package request

import "github.com/final-project-alterra/hospital-management-system-api/features/nurses"

type CreateNurseRequest struct {
	CreatedBy int

	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required"`
	Age      int    `json:"age"`
	ImageUrl string `json:"imageUrl"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Gender   string `json:"gender"`
}

func (c CreateNurseRequest) ToCore() nurses.NurseCore {
	return nurses.NurseCore{
		CreatedBy: c.CreatedBy,
		Email:     c.Email,
		Password:  c.Password,
		Name:      c.Name,
		Age:       c.Age,
		ImageUrl:  c.ImageUrl,
		Phone:     c.Phone,
		Address:   c.Address,
		Gender:    c.Gender,
	}
}
