package response

import (
	"time"

	"github.com/final-project-alterra/hospital-management-system-api/features/nurses"
)

type NurseResponse struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	BirthDate string    `json:"birthDate"`
	ImageUrl  string    `json:"imageUrl"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	Gender    string    `json:"gender"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func DetailNurse(n nurses.NurseCore) NurseResponse {
	return NurseResponse{
		ID:        n.ID,
		Email:     n.Email,
		Name:      n.Name,
		BirthDate: n.BirthDate,
		ImageUrl:  n.ImageUrl,
		Phone:     n.Phone,
		Address:   n.Address,
		Gender:    n.Gender,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
}

func ListNurses(n []nurses.NurseCore) []NurseResponse {
	result := make([]NurseResponse, len(n))
	for i := range n {
		result[i] = DetailNurse(n[i])
	}
	return result
}
