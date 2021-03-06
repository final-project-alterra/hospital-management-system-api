package response

import (
	"fmt"
	"time"

	"github.com/final-project-alterra/hospital-management-system-api/config"
	"github.com/final-project-alterra/hospital-management-system-api/features/admins"
)

type AdminResponse struct {
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

func DetailAdmin(a admins.AdminCore) AdminResponse {
	imageUrl := ""
	if a.ImageUrl != "" {
		imageUrl = fmt.Sprintf("%s/static/%s", config.ENV.DOMAIN, a.ImageUrl)
	}

	return AdminResponse{
		ID:        a.ID,
		Email:     a.Email,
		Name:      a.Name,
		BirthDate: a.BirthDate,
		ImageUrl:  imageUrl,
		Phone:     a.Phone,
		Address:   a.Address,
		Gender:    a.Gender,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func ListAdmin(a []admins.AdminCore) []AdminResponse {
	result := make([]AdminResponse, len(a))
	for i, v := range a {
		result[i] = DetailAdmin(v)
	}
	return result
}
