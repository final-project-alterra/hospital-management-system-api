package presentation

import (
	"net/http"
	"strconv"

	"github.com/final-project-alterra/hospital-management-system-api/errors"
	"github.com/final-project-alterra/hospital-management-system-api/features/admins"
	"github.com/final-project-alterra/hospital-management-system-api/features/admins/presentation/request"
	"github.com/final-project-alterra/hospital-management-system-api/features/admins/presentation/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AdminPresentation struct {
	business admins.IBusiness
	validate *validator.Validate
}

func NewAdminPresentation(business admins.IBusiness) *AdminPresentation {
	validate := validator.New()
	_ = validate.RegisterValidation("ValidateCreateAdminBirthDate", request.ValidateCreateAdminBirthDate)
	_ = validate.RegisterValidation("ValidateEditAdminBirthDate", request.ValidateEditAdminBirthDate)

	return &AdminPresentation{
		business: business,
		validate: validate,
	}
}

func (ap *AdminPresentation) GetAdmins(c echo.Context) error {
	status := http.StatusOK
	message := "Success retrieving admins"
	const op errors.Op = "admins.presentation.GetAdmins"

	data, err := ap.business.FindAdmins()
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}

	return response.Success(c, status, message, response.ListAdmin(data))
}

func (ap *AdminPresentation) GetDetailAdmin(c echo.Context) error {
	status := http.StatusOK
	message := "Success retrieving admin detail"
	const op errors.Op = "admins.presentation.GetDetailAdmin"
	var errMessage errors.ErrClientMessage

	id, err := strconv.Atoi(c.Param("adminId"))
	if err != nil {
		errMessage = "Invalid admin id"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindBadRequest))
	}

	data, err := ap.business.FindAdminById(id)
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}

	return response.Success(c, status, message, response.DetailAdmin(data))
}

func (ap *AdminPresentation) PostCreateAdmin(c echo.Context) error {
	status := http.StatusCreated
	message := "Admin created"
	const op errors.Op = "admins.presentation.PostCreateAdmin"
	var errMessage errors.ErrClientMessage

	creator, ok := c.Get("userId").(int)
	if !ok {
		err := errors.New("Invalid admin creator id")
		errMessage = "Invalid token of admin id"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindBadRequest))
	}

	newAdmin := request.CreateAdminRequest{}
	newAdmin.CreatedBy = creator
	err := c.Bind(&newAdmin)
	if err != nil {
		errMessage = "Error request body"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindBadRequest))
	}

	err = ap.validate.Struct(newAdmin)
	if err != nil {
		errMessage = "Invalid request body. Make sure email is correct and password is min 8 length characters"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindUnprocessable))
	}

	err = ap.business.CreateAdmin(newAdmin.ToAdminCore())
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}

	return response.Success(c, status, message, nil)
}

func (ap *AdminPresentation) PutEditAdmin(c echo.Context) error {
	status := http.StatusOK
	message := "Admin updated"
	const op errors.Op = "admins.presentation.PutEditAdmin"
	var errMessage errors.ErrClientMessage

	updater, ok := c.Get("userId").(int)
	if !ok || updater < 1 {
		err := errors.New("Invalid admin id token")
		errMessage = "Invalid token of admin id"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindBadRequest))
	}

	editedAdmin := request.EditAdminRequest{}
	editedAdmin.UpdatedBy = updater
	err := c.Bind(&editedAdmin)
	if err != nil {
		errMessage = "Invalid request body"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindBadRequest))
	}

	err = ap.validate.Struct(editedAdmin)
	if err != nil {
		errMessage = "Invalid request body. Make sure edited admin ID is attached"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindUnprocessable))
	}

	err = ap.business.EditAdmin(editedAdmin.ToAdminCore())
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}

	return response.Success(c, status, message, nil)
}

func (ap *AdminPresentation) PutEditAdminPassword(c echo.Context) error {
	status := http.StatusOK
	message := "Admin password updated"
	const op errors.Op = "admins.presentation.PutEditAdminPassword"
	var errMessage errors.ErrClientMessage

	updater, ok := c.Get("userId").(int)
	if !ok || updater < 1 {
		err := errors.New("Invalid admin id token")
		errMessage = "Invalid token of admin id"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindBadRequest))
	}

	req := request.EditAdminPasswordRequest{}
	req.UpdatedBy = updater
	err := c.Bind(&req)
	if err != nil {
		errMessage = "Invalid request body"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindBadRequest))
	}

	err = ap.validate.Struct(req)
	if err != nil {
		errMessage = "Make sure id, old password, new password is correct"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindUnprocessable))
	}

	err = ap.business.EditAdminPassword(req.Id, updater, req.OldPassword, req.NewPassword)
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}

	return response.Success(c, status, message, nil)
}

func (ap *AdminPresentation) DeleteAdmin(c echo.Context) error {
	status := http.StatusOK
	message := "Admin successfully deleted"
	const op errors.Op = "admins.presentation.DeleteAdmin"
	var errMessage errors.ErrClientMessage

	updater, ok := c.Get("userId").(int)
	if !ok || updater < 1 {
		err := errors.New("Invalid admin id token")
		errMessage = "Invalid token of admin id"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindBadRequest))
	}

	id, err := strconv.Atoi(c.Param("adminId"))
	if err != nil {
		errMessage = "Invalid admin id"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindBadRequest))
	}

	err = ap.business.RemoveAdminById(id, updater)
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}

	return response.Success(c, status, message, nil)
}
