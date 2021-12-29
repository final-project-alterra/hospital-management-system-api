package presentation

import (
	"net/http"
	"strconv"

	"github.com/final-project-alterra/hospital-management-system-api/errors"
	"github.com/final-project-alterra/hospital-management-system-api/features/doctors"
	"github.com/final-project-alterra/hospital-management-system-api/features/doctors/presentation/request"
	"github.com/final-project-alterra/hospital-management-system-api/features/doctors/presentation/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type DoctorPresentation struct {
	business doctors.IBusiness
	valitate *validator.Validate
}

func NewDoctorPresentation(business doctors.IBusiness) *DoctorPresentation {
	return &DoctorPresentation{
		business: business,
		valitate: validator.New(),
	}
}

func (dp *DoctorPresentation) GetDoctors(c echo.Context) error {
	status := http.StatusOK
	message := "Success retrieving doctors"
	const op errors.Op = "doctors.presentation.GetDoctors"

	doctors, err := dp.business.FindDoctors()
	if err != nil {
		return response.Error(c, errors.E(op, err))
	}
	return response.Success(c, status, message, response.ListDoctors(doctors))
}

func (dp *DoctorPresentation) GetDetailDoctor(c echo.Context) error {
	status := http.StatusOK
	message := "Success retrieving detail doctor"
	const op errors.Op = "doctors.presentation.GetDetailDoctor"
	var errMessage errors.ErrClientMessage

	updatedBy, err := strconv.Atoi(c.Param("doctorId"))
	if err != nil || updatedBy < 1 {
		errMessage = "Invalid doctor id"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindBadRequest))
	}

	doctor, err := dp.business.FindDoctorById(updatedBy)
	if err != nil {
		return response.Error(c, errors.E(op, err))
	}
	return response.Success(c, status, message, response.DetailDoctor(doctor))
}

func (dp *DoctorPresentation) PostDoctor(c echo.Context) error {
	status := http.StatusCreated
	message := "Success creating doctor"
	const op errors.Op = "doctors.presentation.PostDoctor"
	var errMessage errors.ErrClientMessage

	updatedBy, ok := c.Get("userId").(int)
	if !ok || updatedBy < 1 {
		err := errors.New("Invalid admin id")
		errMessage = "Invalid admin id"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindBadRequest))
	}

	var doctor request.CreateDoctorRequest
	doctor.CreatedBy = updatedBy

	err := c.Bind(&doctor)
	if err != nil {
		errMessage = "Invalid doctor payload request"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindBadRequest))
	}

	err = dp.valitate.Struct(doctor)
	if err != nil {
		errMessage = "Invalid payload request. Makesure all field is filled correctly"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindUnprocessable))
	}

	err = dp.business.CreateDoctor(doctor.ToDoctorCore())
	if err != nil {
		return response.Error(c, errors.E(op, err))
	}

	return response.Success(c, status, message, nil)
}

func (dp DoctorPresentation) PutEditDoctor(c echo.Context) error {
	status := http.StatusOK
	message := "Success updating doctor profile"
	const op errors.Op = "doctors.presentation.PutEditDoctor"
	var errMessage errors.ErrClientMessage

	updatedBy, ok := c.Get("userId").(int)
	if !ok || updatedBy < 1 {
		err := errors.New("Invalid admin id")
		errMessage = "Invalid admin id"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindBadRequest))
	}

	var doctor request.UpdateDoctorRequest
	doctor.UpdatedBy = updatedBy

	err := c.Bind(&doctor)
	if err != nil {
		errMessage = "Invalid doctor payload request"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindBadRequest))
	}

	err = dp.valitate.Struct(doctor)
	if err != nil {
		errMessage = "Invalid payload request. Makesure all field is filled correctly"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindUnprocessable))
	}

	err = dp.business.EditDoctor(doctor.ToDoctorCore())
	if err != nil {
		return response.Error(c, errors.E(op, err))
	}

	return response.Success(c, status, message, nil)
}

func (dp DoctorPresentation) PutEditDoctorPassword(c echo.Context) error {
	status := http.StatusOK
	message := "Success updating doctor password"
	const op errors.Op = "doctors.presentation.PutEditDoctorPassword"
	var errMessage errors.ErrClientMessage

	updatedBy, ok := c.Get("userId").(int)
	if !ok || updatedBy < 1 {
		err := errors.New("Invalid admin id")
		errMessage = "Invalid admin id"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindBadRequest))
	}

	var doctor request.UpdateDoctorPasswordRequest

	err := c.Bind(&doctor)
	if err != nil {
		errMessage = "Invalid doctor payload request"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindBadRequest))
	}

	err = dp.valitate.Struct(doctor)
	if err != nil {
		errMessage = "Invalid. Makesure all field is filled & new password is 8 character long"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindUnprocessable))
	}

	err = dp.business.EditDoctorPassword(doctor.ID, updatedBy, doctor.OldPassword, doctor.NewPassword)
	if err != nil {
		return response.Error(c, errors.E(op, err))
	}

	return response.Success(c, status, message, nil)
}

func (dp *DoctorPresentation) DeleteDoctor(c echo.Context) error {
	status := http.StatusOK
	message := "Doctor deleted successfully"
	const op errors.Op = "doctors.presentation.DeleteDoctor"
	var errMessage errors.ErrClientMessage

	updatedBy, ok := c.Get("userId").(int)
	if !ok || updatedBy < 1 {
		err := errors.New("Invalid admin id")
		errMessage = "Invalid admin id"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindBadRequest))
	}

	doctorId, err := strconv.Atoi(c.Param("doctorId"))
	if err != nil {
		errMessage = "Invalid doctor id"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindBadRequest))
	}

	err = dp.business.RemoveDoctorById(doctorId, updatedBy)
	if err != nil {
		return response.Error(c, errors.E(op, err))
	}
	return response.Success(c, status, message, nil)
}
