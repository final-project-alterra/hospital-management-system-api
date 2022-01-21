package request

import "github.com/final-project-alterra/hospital-management-system-api/features/nurses"

type UpdateNurseRequest struct {
	ID        int `json:"id" validate:"gt=0"`
	UpdatedBy int

	Name      string `json:"name" validate:"required"`
	BirthDate string `json:"birthDate" validate:"required,ValidateBirthDate"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Gender    string `json:"gender" validate:"required,oneof='L' 'P'"`
}

func (c UpdateNurseRequest) ToCore() nurses.NurseCore {
	return nurses.NurseCore{
		ID:        c.ID,
		UpdatedBy: c.UpdatedBy,
		Name:      c.Name,
		BirthDate: c.BirthDate,
		Phone:     c.Phone,
		Address:   c.Address,
		Gender:    c.Gender,
	}
}
