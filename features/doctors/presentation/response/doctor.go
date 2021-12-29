package response

import "github.com/final-project-alterra/hospital-management-system-api/features/doctors"

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

	Name     string `json:"name"`
	Email    string `json:"email"`
	ImageUrl string `json:"image_url"`
	Address  string `json:"address"`
	Age      int    `json:"age"`
	Phone    string `json:"phone"`
	Gender   string `json:"gender"`
}

func DetailDoctor(d doctors.DoctorCore) DoctorResponse {
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

		Name:     d.Name,
		Email:    d.Email,
		ImageUrl: d.ImageUrl,
		Address:  d.Address,
		Age:      d.Age,
		Phone:    d.Phone,
		Gender:   d.Gender,
	}
}

func ListDoctors(d []doctors.DoctorCore) []DoctorResponse {
	result := make([]DoctorResponse, len(d))
	for i := range d {
		result[i] = DetailDoctor(d[i])
	}
	return result
}
