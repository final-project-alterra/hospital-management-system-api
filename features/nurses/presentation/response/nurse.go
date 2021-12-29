package response

import "github.com/final-project-alterra/hospital-management-system-api/features/nurses"

type NurseResponse struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	ImageUrl string `json:"imageUrl"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Gender   string `json:"gender"`
}

func DetailNurse(n nurses.NurseCore) NurseResponse {
	return NurseResponse{
		ID:       n.ID,
		Email:    n.Email,
		Name:     n.Name,
		Age:      n.Age,
		ImageUrl: n.ImageUrl,
		Phone:    n.Phone,
		Address:  n.Address,
		Gender:   n.Gender,
	}
}

func ListNurses(n []nurses.NurseCore) []NurseResponse {
	result := make([]NurseResponse, len(n))
	for i := range n {
		result[i] = DetailNurse(n[i])
	}
	return result
}
