package presentation

import (
	"net/http"
	"strconv"

	"github.com/final-project-alterra/hospital-management-system-api/errors"
	"github.com/final-project-alterra/hospital-management-system-api/features/patients"
	"github.com/final-project-alterra/hospital-management-system-api/features/patients/presentation/request"
	"github.com/final-project-alterra/hospital-management-system-api/features/patients/presentation/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type PatientPresentation struct {
	business patients.IBusiness
	validate *validator.Validate
}

func NewPatientPresentation(b patients.IBusiness) *PatientPresentation {
	return &PatientPresentation{
		business: b,
		validate: validator.New(),
	}
}

func (p *PatientPresentation) GetPatients(c echo.Context) error {
	status := http.StatusOK
	message := "Success retrieving patients"
	const op errors.Op = "patients.presentation.GetPatients"

	patientsData, err := p.business.FindPatients()
	if err != nil {
		return response.Error(c, errors.E(op, err))
	}
	return response.Success(c, status, message, response.ListPatients(patientsData))
}

func (p *PatientPresentation) GetDetailPatient(c echo.Context) error {
	status := http.StatusOK
	message := "Success retrieving detail patient"
	const op errors.Op = "patients.presentation.GetDetailPatient"
	var errMessage errors.ErrClientMessage

	patientId, err := strconv.Atoi(c.Param("patientId"))
	if err != nil {
		errMessage = "Invalid patient id"
		return errors.E(err, op, errMessage, errors.KindBadRequest)
	}

	patientData, err := p.business.FindPatientById(patientId)
	if err != nil {
		return response.Error(c, errors.E(op, err))
	}
	return response.Success(c, status, message, response.DetailPatient(patientData))
}

func (p *PatientPresentation) PostPatient(c echo.Context) error {
	status := http.StatusCreated
	message := "Success creating patient"
	const op errors.Op = "patients.presentation.PostPatient"
	var errMessage errors.ErrClientMessage

	createdBy, ok := c.Get("userId").(int)
	if !ok {
		err := errors.New("Invalid admin id")
		errMessage = "Invalid admin id"
		return errors.E(err, op, errMessage, errors.KindBadRequest)
	}

	patient := request.CreatePatientRequest{CreatedBy: createdBy}
	if err := c.Bind(&patient); err != nil {
		errMessage = "Unable to parse data"
		return errors.E(err, op, errMessage, errors.KindBadRequest)
	}

	if err := p.validate.Struct(&patient); err != nil {
		errMessage = "Invalid data. Makesure all data is filled correctly"
		return errors.E(err, op, errMessage, errors.KindUnprocessable)
	}

	if err := p.business.CreatePatient(patient.ToPatientCore()); err != nil {
		return response.Error(c, errors.E(op, err))
	}
	return response.Success(c, status, message, nil)
}
