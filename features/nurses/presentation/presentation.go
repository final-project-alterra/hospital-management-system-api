package presentation

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/final-project-alterra/hospital-management-system-api/errors"
	"github.com/final-project-alterra/hospital-management-system-api/features/nurses"
	"github.com/final-project-alterra/hospital-management-system-api/features/nurses/presentation/request"
	"github.com/final-project-alterra/hospital-management-system-api/features/nurses/presentation/response"
	"github.com/final-project-alterra/hospital-management-system-api/utils/project"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type NursePresentation struct {
	business nurses.IBusiness
	validate *validator.Validate
}

func NewNursePresentation(business nurses.IBusiness) *NursePresentation {
	validate := validator.New()
	_ = validate.RegisterValidation("ValidateBirthDate", request.ValidateBirthDate)

	return &NursePresentation{
		business: business,
		validate: validate,
	}
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

	nurseId, err := strconv.Atoi(c.Param("nurseId"))
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
	const op errors.Op = "presentation.nurses.PutEditNurse"
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
	const op errors.Op = "presentation.nurses.PutEditNursePassword"
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

func (ap *NursePresentation) PutEditImageProfile(c echo.Context) error {
	status := http.StatusOK
	message := "Image profile updated"
	const op errors.Op = "nurses.presentation.PutEditImageProfile"
	var errMsg errors.ErrClientMessage

	updatedBy := c.Get("userId").(int)
	nuserID, err := strconv.Atoi(c.FormValue("nurseId"))
	if err != nil {
		errMsg = "Unable to parse nurse id"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindBadRequest))
	}

	destDirectory := path.Join(project.GetMainDir(), "files")
	filename, err := ap.allocateFile(c, destDirectory)
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}

	updatedNurse := nurses.NurseCore{ID: nuserID, UpdatedBy: updatedBy, ImageUrl: filename}
	err = ap.business.EditNurseImageProfile(updatedNurse)
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}

	return response.Success(c, status, message, nil)
}

func (ap *NursePresentation) DeleteImageProfile(c echo.Context) error {
	status := http.StatusOK
	message := "Image profile deleted"
	const op errors.Op = "admins.presentation.DeleteImageProfile"
	var errMsg errors.ErrClientMessage

	updatedBy := c.Get("userId").(int)
	nurseID, err := strconv.Atoi(c.Param("nurseId"))
	if err != nil {
		errMsg = "Unable to parse nurse id"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindBadRequest))
	}

	updatedNurse := nurses.NurseCore{ID: nurseID, UpdatedBy: updatedBy, ImageUrl: ""}
	err = ap.business.EditNurseImageProfile(updatedNurse)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, status, message, nil)
}

func (np *NursePresentation) DeleteNurse(c echo.Context) error {
	status := http.StatusOK
	message := "Success deleting nurse"
	const op errors.Op = "presentation.nurses.DeleteNurse"
	var errMessage errors.ErrClientMessage

	updatedBy, ok := c.Get("userId").(int)
	if !ok || updatedBy < 1 {
		err := errors.New("Invalid admin id")
		errMessage = "Invalid admin id"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindBadRequest))
	}

	nurseId, err := strconv.Atoi(c.Param("nurseId"))
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

// Private methods
func (ap *NursePresentation) allocateFile(c echo.Context, destDirectory string) (string, error) {
	const op errors.Op = "nurses.presentation.allocateFile"
	var errMsg errors.ErrClientMessage = "Something went wrong"

	var file *multipart.FileHeader
	var src multipart.File
	var dst *os.File
	var err error

	if err = os.MkdirAll(destDirectory, os.ModePerm); err != nil {
		return "", errors.E(err, op, errMsg, errors.KindServerError)
	}

	if file, err = c.FormFile("image"); err != nil {
		errMsg = "Unable to parse image"
		return "", errors.E(err, op, errMsg, errors.KindBadRequest)
	}

	if src, err = file.Open(); err != nil {
		return "", errors.E(err, op, errMsg, errors.KindServerError)
	}
	defer src.Close()

	filename := fmt.Sprintf("%s-%s", uuid.New().String(), file.Filename)
	filepath := path.Join(destDirectory, filename)
	if dst, err = os.Create(filepath); err != nil {
		return "", errors.E(err, op, errMsg, errors.KindServerError)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", errors.E(err, op, errMsg, errors.KindServerError)
	}

	return filename, err
}
