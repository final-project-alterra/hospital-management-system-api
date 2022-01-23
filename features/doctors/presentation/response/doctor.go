package response

import (
	"fmt"
	"time"

	"github.com/final-project-alterra/hospital-management-system-api/config"
	"github.com/final-project-alterra/hospital-management-system-api/features/doctors"
)

type DoctorSpecialityResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type DoctorRoomResponse struct {
	ID    int    `json:"id"`
	Floor string `json:"floor"`
	Code  string `json:"code"`
}

type DoctorResponse struct {
	ID int `json:"id"`

	Speciality DoctorSpecialityResponse `json:"speciality"`
	Room       DoctorRoomResponse       `json:"room"`

	Name      string    `json:"name"`
	Email     string    `json:"email"`
	ImageUrl  string    `json:"imageUrl"`
	Address   string    `json:"address"`
	BirthDate string    `json:"birthDate"`
	Phone     string    `json:"phone"`
	Gender    string    `json:"gender"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func DetailDoctor(d doctors.DoctorCore) DoctorResponse {
	imageUrl := ""
	if d.ImageUrl != "" {
		imageUrl = fmt.Sprintf("%s/static/%s", config.ENV.DOMAIN, d.ImageUrl)
	}

	return DoctorResponse{
		ID: d.ID,

		Speciality: DoctorSpecialityResponse{
			ID:   d.Speciality.ID,
			Name: d.Speciality.Name,
		},
		Room: DoctorRoomResponse{
			ID:    d.Room.ID,
			Floor: d.Room.Floor,
			Code:  d.Room.Code,
		},

		Name:      d.Name,
		Email:     d.Email,
		ImageUrl:  imageUrl,
		Address:   d.Address,
		BirthDate: d.BirthDate,
		Phone:     d.Phone,
		Gender:    d.Gender,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

func ListDoctors(d []doctors.DoctorCore) []DoctorResponse {
	result := make([]DoctorResponse, len(d))
	for i := range d {
		result[i] = DetailDoctor(d[i])
	}
	return result
}
