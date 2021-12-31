package request

type UpdateDoctorPasswordRequest struct {
	ID          int    `json:"id" validate:"gt=0"`
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=8"`
}
