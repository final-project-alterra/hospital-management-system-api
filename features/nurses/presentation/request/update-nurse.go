package request

import "github.com/final-project-alterra/hospital-management-system-api/features/nurses"

type UpdateNurseRequest struct {
	ID        int `json:"id" validate:"gt=0"`
	UpdatedBy int

	Name     string `json:"name" validate:"required"`
	Age      int    `json:"age"`
	ImageUrl string `json:"imageUrl"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Gender   string `json:"gender"`
}

func (c UpdateNurseRequest) ToCore() nurses.NurseCore {
	return nurses.NurseCore{
		ID:        c.ID,
		UpdatedBy: c.UpdatedBy,
		Name:      c.Name,
		Age:       c.Age,
		ImageUrl:  c.ImageUrl,
		Phone:     c.Phone,
		Address:   c.Address,
		Gender:    c.Gender,
	}
}
