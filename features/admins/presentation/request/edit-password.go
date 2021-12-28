package request

type EditAdminPasswordRequest struct {
	Id          int `json:"id" validate:"required"`
	UpdatedBy   int
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=8"`
}
