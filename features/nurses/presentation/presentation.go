package presentation

import (
	"net/http"
	"strconv"

	"github.com/final-project-alterra/hospital-management-system-api/errors"
	"github.com/final-project-alterra/hospital-management-system-api/features/nurses"
	"github.com/final-project-alterra/hospital-management-system-api/features/nurses/presentation/request"
	"github.com/final-project-alterra/hospital-management-system-api/features/nurses/presentation/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type NursePresentation struct {
	business nurses.IBusiness
	validate *validator.Validate
}

func NewNursePresentation(business nurses.IBusiness) *NursePresentation {
	return &NursePresentation{business, validator.New()}
}

func (np *NursePresentation) GetNurses(c echo.Context) error {
	status := http.StatusOK
	message := "Success retrieving nurses data"
	const op errors.Op = "presentation.nurses.GetNurses"

	nursesData, err := np.business.FindNurses()
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}
	return response.Success(c, status, message, response.ListNurses(nursesData))
}

func (np *NursePresentation) GetDetailNurse(c echo.Context) error {
	status := http.StatusOK
	message := "Success retrieving nurses data"
	const op errors.Op = "presentation.nurses.GetDetailNurse"
	var errMessage errors.ErrClientMessage

	nurseId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errMessage = "Invalid nurse id"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindBadRequest))
	}

	nurseData, err := np.business.FindNurseById(nurseId)
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}
	return response.Success(c, status, message, response.DetailNurse(nurseData))
}

func (np *NursePresentation) PostNurse(c echo.Context) error {
	status := http.StatusCreated
	message := "Success creating nurse"
	const op errors.Op = "presentation.nurses.PostNurse"
	var errMessage errors.ErrClientMessage

	createdBy, ok := c.Get("userId").(int)
	if !ok || createdBy < 1 {
		err := errors.New("Invalid admin id")
		errMessage = "Invalid admin id"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindBadRequest))
	}

	nurse := request.CreateNurseRequest{CreatedBy: createdBy}
	if err := c.Bind(&nurse); err != nil {
		errMessage = "Unable to parse nurse request payload"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindBadRequest))
	}

	if err := np.validate.Struct(nurse); err != nil {
		errMessage = "Invalid nurse data. Makesure all required fields are filled correctly"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindUnprocessable))
	}

	if err := np.business.CreateNurse(nurse.ToCore()); err != nil {
		return response.Error(c, errors.E(err, op))
	}

	return response.Success(c, status, message, nil)
}

func (np *NursePresentation) PutEditNurse(c echo.Context) error {
	status := http.StatusOK
	message := "Success updating nurse profile data"
	const op errors.Op = "presentation.nurses.PostNurse"
	var errMessage errors.ErrClientMessage

	updatedBy, ok := c.Get("userId").(int)
	if !ok || updatedBy < 1 {
		err := errors.New("Invalid admin id")
		errMessage = "Invalid admin id"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindBadRequest))
	}

	nurse := request.UpdateNurseRequest{UpdatedBy: updatedBy}
	if err := c.Bind(&nurse); err != nil {
		errMessage = "Unable to parse nurse request payload"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindBadRequest))
	}

	if err := np.validate.Struct(nurse); err != nil {
		errMessage = "Invalid nurse data. Makesure all required fields are filled correctly"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindUnprocessable))
	}

	if err := np.business.EditNurse(nurse.ToCore()); err != nil {
		return response.Error(c, errors.E(err, op))
	}

	return response.Success(c, status, message, nil)
}

func (np *NursePresentation) PutEditNursePassword(c echo.Context) error {
	status := http.StatusOK
	message := "Success editing nurse password"
	const op errors.Op = "presentation.nurses.PostNurse"
	var errMessage errors.ErrClientMessage

	updatedBy, ok := c.Get("userId").(int)
	if !ok || updatedBy < 1 {
		err := errors.New("Invalid admin id")
		errMessage = "Invalid admin id"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindBadRequest))
	}

	nurse := request.UpdateNursePasswordRequest{}
	if err := c.Bind(&nurse); err != nil {
		errMessage = "Unable to parse nurse request payload"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindBadRequest))
	}

	if err := np.validate.Struct(nurse); err != nil {
		errMessage = "Invalid nurse data. Makesure all fields are filled correctly and new password is min 8 characters"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindUnprocessable))
	}

	err := np.business.EditNursePassword(nurse.ID, updatedBy, nurse.OldPassword, nurse.NewPassword)
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}

	return response.Success(c, status, message, nil)
}

func (np *NursePresentation) DeleteNurse(c echo.Context) error {
	status := http.StatusOK
	message := "Success deleting nurse"
	const op errors.Op = "presentation.nurses.PostNurse"
	var errMessage errors.ErrClientMessage

	updatedBy, ok := c.Get("userId").(int)
	if !ok || updatedBy < 1 {
		err := errors.New("Invalid admin id")
		errMessage = "Invalid admin id"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindBadRequest))
	}

	nurseId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		err := errors.New("Invalid nurse id")
		errMessage = "Invalid nurse id"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindBadRequest))
	}

	err = np.business.RemoveNurseById(nurseId, updatedBy)
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}

	return response.Success(c, status, message, nil)
}
