package request

import "github.com/final-project-alterra/hospital-management-system-api/features/admins"

type EditAdminRequest struct {
	ID        int `json:"id" validate:"required"`
	UpdatedBy int
	Name      string `json:"name"`
	Age       int    `json:"age"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Gender    string `json:"gender"`
}

func (r EditAdminRequest) ToAdminCore() admins.AdminCore {
	return admins.AdminCore{
		ID:        r.ID,
		UpdatedBy: r.UpdatedBy,
		Name:      r.Name,
		Age:       r.Age,
		Phone:     r.Phone,
		Address:   r.Address,
		Gender:    r.Gender,
	}
}
