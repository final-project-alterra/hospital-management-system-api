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
	"github.com/final-project-alterra/hospital-management-system-api/features/doctors"
	"github.com/final-project-alterra/hospital-management-system-api/features/doctors/presentation/request"
	"github.com/final-project-alterra/hospital-management-system-api/features/doctors/presentation/response"
	"github.com/final-project-alterra/hospital-management-system-api/utils/project"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type DoctorPresentation struct {
	business doctors.IBusiness
	valitate *validator.Validate
}

func NewDoctorPresentation(business doctors.IBusiness) *DoctorPresentation {
	validate := validator.New()
	_ = validate.RegisterValidation("ValidateBirthDate", request.ValidateBirthDate)

	return &DoctorPresentation{
		business: business,
		valitate: validate,
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

	createdBy, ok := c.Get("userId").(int)
	if !ok || createdBy < 1 {
		err := errors.New("Invalid admin id")
		errMessage = "Invalid admin id"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindBadRequest))
	}

	var doctor request.CreateDoctorRequest
	doctor.CreatedBy = createdBy

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

func (ap *DoctorPresentation) PutEditImageProfile(c echo.Context) error {
	status := http.StatusOK
	message := "Image profile updated"
	const op errors.Op = "doctors.presentation.PutEditImageProfile"
	var errMsg errors.ErrClientMessage

	updatedBy := c.Get("userId").(int)
	doctorID, err := strconv.Atoi(c.FormValue("doctorId"))
	if err != nil {
		errMsg = "Unable to parse doctor id"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindBadRequest))
	}

	destDirectory := path.Join(project.GetMainDir(), "files")
	filename, err := ap.allocateFile(c, destDirectory)
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}

	updatedDoctor := doctors.DoctorCore{ID: doctorID, UpdatedBy: updatedBy, ImageUrl: filename}
	err = ap.business.EditDoctorImageProfile(updatedDoctor)
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}

	return response.Success(c, status, message, nil)
}

func (ap *DoctorPresentation) DeleteImageProfile(c echo.Context) error {
	status := http.StatusOK
	message := "Image profile deleted"
	const op errors.Op = "doctors.presentation.DeleteImageProfile"
	var errMsg errors.ErrClientMessage

	updatedBy := c.Get("userId").(int)
	doctorID, err := strconv.Atoi(c.Param("doctorId"))
	if err != nil {
		errMsg = "Unable to parse doctor id"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindBadRequest))
	}

	updatedDoctor := doctors.DoctorCore{ID: doctorID, UpdatedBy: updatedBy, ImageUrl: ""}
	err = ap.business.EditDoctorImageProfile(updatedDoctor)
	if err != nil {
		return response.Error(c, errors.E(err, op))
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

func (dp *DoctorPresentation) GetSpecialities(c echo.Context) error {
	status := http.StatusOK
	message := "Success retrieving specialities"
	const op errors.Op = "doctors.presentation.GetSpecialities"

	specialities, err := dp.business.FindSpecialities()
	if err != nil {
		return response.Error(c, errors.E(op, err))
	}

	return response.Success(c, status, message, response.ListSpecialities(specialities))
}

func (dp *DoctorPresentation) GetDetailSpeciality(c echo.Context) error {
	status := http.StatusOK
	message := "Success retrieving detail speciality"
	const op errors.Op = "doctors.presentation.GetDetailSpeciality"
	var errMessage errors.ErrClientMessage

	specialityId, err := strconv.Atoi(c.Param("specialityId"))
	if err != nil {
		errMessage = "Invalid id speciality"
		return response.Error(c, errors.E(op, err, errMessage, errors.KindBadRequest))
	}

	speciality, err := dp.business.FindSpecialityById(specialityId)
	if err != nil {
		return response.Error(c, errors.E(op, err))
	}

	return response.Success(c, status, message, response.DetailSpeciality(speciality))
}

func (dp *DoctorPresentation) PostSpeciality(c echo.Context) error {
	status := http.StatusCreated
	message := "Success creating speciality"
	const op errors.Op = "doctors.presentation.PostSpeciality"
	var errMessage errors.ErrClientMessage

	speciality := request.CreateSpecialityRequest{}
	err := c.Bind(&speciality)
	if err != nil {
		errMessage = "Invalid payload request"
		return response.Error(c, errors.E(op, err, errMessage, errors.KindBadRequest))
	}

	err = dp.valitate.Struct(speciality)
	if err != nil {
		errMessage = "Invalid. Makesure speciality name is filled correctly"
		return response.Error(c, errors.E(op, err, errMessage, errors.KindUnprocessable))
	}

	err = dp.business.CreateSpeciality(speciality.ToSpecialityCore())
	if err != nil {
		return response.Error(c, errors.E(op, err))
	}

	return response.Success(c, status, message, nil)
}

func (dp *DoctorPresentation) PutEditSpeciality(c echo.Context) error {
	status := http.StatusOK
	message := "Success updating speciality"
	const op errors.Op = "doctors.presentation.PutEditSpeciality"
	var errMessage errors.ErrClientMessage

	speciality := request.UpdateSpecialityRequest{}
	err := c.Bind(&speciality)
	if err != nil {
		errMessage = "Unable to parse payload request"
		return response.Error(c, errors.E(op, err, errMessage, errors.KindBadRequest))
	}

	err = dp.valitate.Struct(speciality)
	if err != nil {
		errMessage = "Invalid. Makesure all field is filled correctly"
		return response.Error(c, errors.E(op, err, errMessage, errors.KindUnprocessable))
	}

	err = dp.business.EditSpeciality(speciality.ToSpecialityCore())
	if err != nil {
		return response.Error(c, errors.E(op, err))
	}

	return response.Success(c, status, message, nil)
}

func (dp *DoctorPresentation) DeleteSpeciality(c echo.Context) error {
	status := http.StatusOK
	message := "Success deleting speciality"
	const op errors.Op = "doctors.presentation.DeleteSpeciality"
	var errMessage errors.ErrClientMessage

	specialityId, err := strconv.Atoi(c.Param("specialityId"))
	if err != nil {
		errMessage = "Invalid id speciality"
		return response.Error(c, errors.E(op, err, errMessage, errors.KindBadRequest))
	}

	err = dp.business.RemoveSpeciality(specialityId)
	if err != nil {
		return response.Error(c, errors.E(op, err))
	}

	return response.Success(c, status, message, nil)
}

func (dp *DoctorPresentation) GetRooms(c echo.Context) error {
	status := http.StatusOK
	message := "Success retrieving rooms"
	const op errors.Op = "doctors.presentation.GetRooms"

	rooms, err := dp.business.FindRooms()
	if err != nil {
		return response.Error(c, errors.E(op, err))
	}

	return response.Success(c, status, message, response.ListRooms(rooms))
}

func (dp *DoctorPresentation) PostRoom(c echo.Context) error {
	status := http.StatusCreated
	message := "Success creating rooms"
	const op errors.Op = "doctors.presentation.PostRoom"
	var errMessage errors.ErrClientMessage

	room := request.CreateRoomRequest{}
	err := c.Bind(&room)
	if err != nil {
		errMessage = "Unable to parse payload request"
		return response.Error(c, errors.E(op, err, errMessage, errors.KindBadRequest))
	}

	err = dp.valitate.Struct(room)
	if err != nil {
		errMessage = "Invalid. Makesure all field is filled correctly"
		return response.Error(c, errors.E(op, err, errMessage, errors.KindUnprocessable))
	}

	err = dp.business.CreateRoom(room.ToRoomCore())
	if err != nil {
		return response.Error(c, errors.E(op, err))
	}
	return response.Success(c, status, message, nil)
}

func (dp *DoctorPresentation) PutEditRoom(c echo.Context) error {
	status := http.StatusOK
	message := "Success updating room"
	const op errors.Op = "doctors.presentation.PutEditRoom"
	var errMessage errors.ErrClientMessage

	room := request.UpdateRoomRequest{}
	err := c.Bind(&room)
	if err != nil {
		errMessage = "Unable to parse payload request"
		return response.Error(c, errors.E(op, err, errMessage, errors.KindBadRequest))
	}

	err = dp.valitate.Struct(room)
	if err != nil {
		errMessage = "Invalid. Makesure all field is filled correctly"
		return response.Error(c, errors.E(op, err, errMessage, errors.KindUnprocessable))
	}

	err = dp.business.EditRoom(room.ToRoomCore())
	if err != nil {
		return response.Error(c, errors.E(op, err))
	}

	return response.Success(c, status, message, nil)
}

func (dp *DoctorPresentation) DeleteRoom(c echo.Context) error {
	status := http.StatusOK
	message := "Success deleting room"
	const op errors.Op = "doctors.presentation.DeleteRoom"
	var errMessage errors.ErrClientMessage

	roomId, err := strconv.Atoi(c.Param("roomId"))
	if err != nil || roomId < 1 {
		errMessage = "Invalid id room"
		return response.Error(c, errors.E(op, err, errMessage, errors.KindBadRequest))
	}

	err = dp.business.RemoveRoomById(roomId)
	if err != nil {
		return response.Error(c, errors.E(op, err))
	}
	return response.Success(c, status, message, nil)
}

// Private methods
func (ap *DoctorPresentation) allocateFile(c echo.Context, destDirectory string) (string, error) {
	const op errors.Op = "doctors.presentation.allocateFile"
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
