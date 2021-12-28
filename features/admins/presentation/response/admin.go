package response

import "github.com/final-project-alterra/hospital-management-system-api/features/admins"

type AdminResponse struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	ImageUrl string `json:"imageUrl"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Gender   string `json:"gender"`
}

func DetailAdmin(a admins.AdminCore) AdminResponse {
	return AdminResponse{
		ID:       a.ID,
		Email:    a.Email,
		Name:     a.Name,
		Age:      a.Age,
		ImageUrl: a.ImageUrl,
		Phone:    a.Phone,
		Address:  a.Address,
		Gender:   a.Gender,
	}
}

func ListAdmin(a []admins.AdminCore) []AdminResponse {
	result := make([]AdminResponse, len(a))
	for i, v := range a {
		result[i] = DetailAdmin(v)
	}
	return result
}
