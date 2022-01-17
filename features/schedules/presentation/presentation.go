package presentation

import (
	"net/http"
	"strconv"

	"github.com/final-project-alterra/hospital-management-system-api/errors"
	"github.com/final-project-alterra/hospital-management-system-api/features/schedules"
	"github.com/final-project-alterra/hospital-management-system-api/features/schedules/presentation/request"
	"github.com/final-project-alterra/hospital-management-system-api/features/schedules/presentation/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type SchedulePresentation struct {
	business schedules.IBusiness
	validate *validator.Validate
}

func NewSchedulePresentation(business schedules.IBusiness) *SchedulePresentation {
	validate := validator.New()
	_ = validate.RegisterValidation("ValidateQueryDate", request.ValidateQueryDate)
	_ = validate.RegisterValidation("ValidateQueryLimit", request.ValidateQueryLimit)
	_ = validate.RegisterValidation("ValidateCreateScheduleDate", request.ValidateCreateScheduleDate)
	_ = validate.RegisterValidation("ValidateCreateScheduleTime", request.ValidateCreateScheduleTime)
	_ = validate.RegisterValidation("ValidateUpdateScheduleTime", request.ValidateUpdateScheduleTime)
	_ = validate.RegisterValidation("ValidateUpdateScheduleDate", request.ValidateUpdateScheduleDate)

	return &SchedulePresentation{
		business: business,
		validate: validate,
	}
}

/* Work schedules */
func (p *SchedulePresentation) GetWorkSchedules(c echo.Context) error {
	const op errors.Op = "schedules.presentation.GetWorkSchedules"
	var errMsg errors.ErrClientMessage

	code := http.StatusOK
	message := "Successfully get work schedules"

	query := request.NewQueryParamsRequest()
	if err := c.Bind(&query); err != nil {
		errMsg = "Unable to parse query params"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindBadRequest))
	}

	if err := p.validate.Struct(query); err != nil {
		errMsg = "Invalid query. Makesure date in the format of YYYY-MM-DD"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindBadRequest))
	}

	schedulesData, err := p.business.FindWorkSchedules(query.ToScheduleQuery())
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}

	return response.Success(c, code, message, response.ListWorkSchedule(schedulesData))
}

func (p *SchedulePresentation) GetDoctorWorkSchedules(c echo.Context) error {
	const op errors.Op = "schedules.presentation.GetDoctorWorkSchedules"
	var errMsg errors.ErrClientMessage

	code := http.StatusOK
	message := "Successfully get work schedules"

	query := request.NewQueryParamsRequest()
	if err := c.Bind(&query); err != nil {
		errMsg = "Unable to parse query params"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindBadRequest))
	}

	doctorID, err := strconv.Atoi(c.Param("doctorId"))
	if err != nil {
		errMsg = "Invalid doctor id"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindBadRequest))
	}

	if err := p.validate.Struct(query); err != nil {
		errMsg = "Invalid query. Makesure date in the format of YYYY-MM-DD"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindUnprocessable))
	}

	schedulesData, err := p.business.FindDoctorWorkSchedules(doctorID, query.ToScheduleQuery())
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}

	return response.Success(c, code, message, response.ListDoctorSchedule(schedulesData))
}

func (p *SchedulePresentation) GetNurseWorkSchedules(c echo.Context) error {
	const op errors.Op = "schedules.presentation.GetNurseWorkSchedules"
	var errMsg errors.ErrClientMessage

	code := http.StatusOK
	message := "Successfully get work schedules"

	query := request.NewQueryParamsRequest()
	if err := c.Bind(&query); err != nil {
		errMsg = "Unable to parse query params"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindBadRequest))
	}

	nurseID, err := strconv.Atoi(c.Param("nurseId"))
	if err != nil {
		errMsg = "Invalid nurse id"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindBadRequest))
	}

	if err := p.validate.Struct(query); err != nil {
		errMsg = "Invalid query. Makesure date in the format of YYYY-MM-DD"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindUnprocessable))
	}

	schedulesData, err := p.business.FindNurseWorkSchedules(nurseID, query.ToScheduleQuery())
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}

	return response.Success(c, code, message, response.ListNurseSchedule(schedulesData))
}

func (p *SchedulePresentation) PostWorkSchedules(c echo.Context) error {
	const op errors.Op = "schedules.presentation.PostWorkSchedules"
	var errMsg errors.ErrClientMessage

	code := http.StatusCreated
	message := "Successfully creating schedules"

	schedule := request.CreateWorkScheduleRequest{}
	if err := c.Bind(&schedule); err != nil {
		errMsg = "Unable to parse request body"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindBadRequest))
	}

	if err := p.validate.Struct(schedule); err != nil {
		errMsg = "Invalid request. Makesure all fields are filled correctly"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindUnprocessable))
	}

	workSchedule := schedules.WorkScheduleCore{}
	workSchedule.Doctor.ID = schedule.DoctorID
	workSchedule.Nurse.ID = schedule.NurseID
	workSchedule.StartTime = schedule.StartTime
	workSchedule.EndTime = schedule.EndTime

	query := schedules.ScheduleQuery{}
	query.Repeat = schedule.Repeat
	query.StartDate = schedule.StartDate
	query.EndDate = schedule.EndDate

	err := p.business.CreateWorkSchedule(workSchedule, query)
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}

	return response.Success(c, code, message, nil)
}

func (p *SchedulePresentation) PutEditWorkSchedule(c echo.Context) error {
	const op errors.Op = "schedules.presentation.PutEditWorkSchedule"
	var errMsg errors.ErrClientMessage

	code := http.StatusOK
	message := "Successfully updating schedule"

	updatedSchedule := request.UpdateWorkScheduleRequest{}
	if err := c.Bind(&updatedSchedule); err != nil {
		errMsg = "Unable to parse request body"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindBadRequest))
	}

	if err := p.validate.Struct(updatedSchedule); err != nil {
		errMsg = "Invalid request. Makesure all fields are filled correctly"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindUnprocessable))
	}

	err := p.business.EditWorkSchedule(updatedSchedule.ToWorkScheduleCore())
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}

	return response.Success(c, code, message, nil)
}

func (p *SchedulePresentation) DeleteWorkSchedule(c echo.Context) error {
	const op errors.Op = "schedules.presentation.DeleteWorkSchedule"
	var errMsg errors.ErrClientMessage

	code := http.StatusOK
	message := "Successfully deleting schedule"

	scheduleID, err := strconv.Atoi(c.Param("workScheduleId"))
	if err != nil {
		errMsg = "Invalid work schedule id"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindBadRequest))
	}

	err = p.business.RemoveWorkScheduleById(scheduleID)
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}

	return response.Success(c, code, message, nil)
}

/* Outpatients */
func (p *SchedulePresentation) GetOutpatients(c echo.Context) error {
	const op errors.Op = "schedules.presentation.GetOutpatients"
	var errMsg errors.ErrClientMessage

	code := http.StatusOK
	message := "Successfully retrieving outpatients"

	query := request.NewQueryParamsRequest()
	if err := c.Bind(&query); err != nil {
		errMsg = "Unable to parse query params"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindBadRequest))
	}

	if err := p.validate.Struct(query); err != nil {
		errMsg = "Invalid query. Make sure all query is valid"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindUnprocessable))
	}

	outpatients, err := p.business.FindOutpatients(query.ToScheduleQuery())
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}

	return response.Success(c, code, message, response.ListOutpatients(outpatients))
}

func (p *SchedulePresentation) GetPatientOutpatients(c echo.Context) error {
	const op errors.Op = "schedules.presentation.GetPatientOutpatients"
	var errMsg errors.ErrClientMessage

	code := http.StatusOK
	message := "Successfully retrieving patient outpatients"

	patientID, err := strconv.Atoi(c.Param("patientId"))
	if err != nil {
		errMsg = "Invalid patient id"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindBadRequest))
	}

	query := request.NewQueryParamsRequest()
	if err := c.Bind(&query); err != nil {
		errMsg = "Unable to parse query params"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindBadRequest))
	}

	if err := p.validate.Struct(query); err != nil {
		errMsg = "Invalid query. Make sure all query is valid"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindUnprocessable))
	}

	outpatients, err := p.business.FindOutpatientsByPatientId(patientID, query.ToScheduleQuery())
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}

	return response.Success(c, code, message, response.ListPatientOutpatients(outpatients))
}

func (p *SchedulePresentation) GetWorkScheduleOutpatients(c echo.Context) error {
	const op errors.Op = "schedules.presentation.GetWorkScheduleOutpatients"
	var errMsg errors.ErrClientMessage

	code := http.StatusOK
	message := "Successfully retrieving work schedule outpatients"

	workScheduleID, err := strconv.Atoi(c.Param("workScheduleId"))
	if err != nil {
		errMsg = "Invalid work schedule id"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindBadRequest))
	}

	workSchedule, err := p.business.FindOutpatientsByWorkScheduleId(workScheduleID)
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}

	return response.Success(c, code, message, response.WorkScheduleOutpatient(workSchedule))
}

func (p *SchedulePresentation) GetDetailOutpatient(c echo.Context) error {
	const op errors.Op = "schedules.presentation.GetWorkScheduleOutpatients"
	var errMsg errors.ErrClientMessage

	code := http.StatusOK
	message := "Successfully retrieving work schedule outpatients"

	outpatientID, err := strconv.Atoi(c.Param("outpatientId"))
	if err != nil {
		errMsg = "Invalid outpatient id"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindBadRequest))
	}

	outpatient, err := p.business.FindOutpatientById(outpatientID)
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}

	return response.Success(c, code, message, response.OutpatientDetail(outpatient))
}

func (p *SchedulePresentation) PostOutpatient(c echo.Context) error {
	const op errors.Op = "schedules.presentation.PostOutpatient"
	var errMsg errors.ErrClientMessage

	code := http.StatusCreated
	message := "Successfully creating outpatient"

	outpatient := request.CreateOutpatientRequest{}
	if err := c.Bind(&outpatient); err != nil {
		errMsg = "Unable to parse request body"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindBadRequest))
	}

	if err := p.validate.Struct(outpatient); err != nil {
		errMsg = "Invalid request. Make sure all fields are filled correctly"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindUnprocessable))
	}

	err := p.business.CreateOutpatient(outpatient.ToOutpatientCore())
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}

	return response.Success(c, code, message, nil)
}

func (p *SchedulePresentation) PutEditOutpatient(c echo.Context) error {
	const op errors.Op = "schedules.presentation.PutEditOutpatient"
	var errMsg errors.ErrClientMessage

	code := http.StatusOK
	message := "Successfully updating outpatient"

	outpatient := request.UpdateOutpatientRequest{}
	if err := c.Bind(&outpatient); err != nil {
		errMsg = "Unable to parse request body"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindBadRequest))
	}

	if err := p.validate.Struct(outpatient); err != nil {
		errMsg = "Invalid request. Make sure all fields are filled correctly"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindUnprocessable))
	}

	err := p.business.EditOutpatient(outpatient.ToOutpatientCore())
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}

	return response.Success(c, code, message, nil)
}

func (p *SchedulePresentation) PutExamineOutpatient(c echo.Context) error {
	const op errors.Op = "schedules.presentation.PutExamineOutpatient"
	var errMsg errors.ErrClientMessage

	code := http.StatusOK
	message := "Successfully examined outpatient"

	userID := c.Get("userId").(int)
	role := c.Get("role").(string)
	outpatient := request.ExamineOutpatientRequest{}

	if err := c.Bind(&outpatient); err != nil {
		errMsg = "Unable to parse request body"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindBadRequest))
	}

	if err := p.validate.Struct(outpatient); err != nil {
		errMsg = "Invalid request. Make sure id is correctly"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindUnprocessable))
	}

	err := p.business.ExamineOutpatient(outpatient.ID, userID, role)
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}

	return response.Success(c, code, message, nil)
}

func (p *SchedulePresentation) PutCancelOutpatient(c echo.Context) error {
	const op errors.Op = "schedules.presentation.PutCancelOutpatient"
	var errMsg errors.ErrClientMessage

	code := http.StatusOK
	message := "Successfully canceled outpatient"

	userID := c.Get("userId").(int)
	role := c.Get("role").(string)
	outpatient := request.CancelOutpatientRequest{}

	if err := c.Bind(&outpatient); err != nil {
		errMsg = "Unable to parse request body"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindBadRequest))
	}

	if err := p.validate.Struct(outpatient); err != nil {
		errMsg = "Invalid request. Make sure all fields are filled correctly"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindUnprocessable))
	}

	err := p.business.CancelOutpatient(outpatient.ID, userID, role)
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}

	return response.Success(c, code, message, nil)
}

func (p *SchedulePresentation) PutFinishOutpatient(c echo.Context) error {
	const op errors.Op = "schedules.presentation.PutFinishOutpatient"
	var errMsg errors.ErrClientMessage

	code := http.StatusOK
	message := "Successfully finishing outpatient"

	userID := c.Get("userId").(int)
	role := c.Get("role").(string)
	outpatient := request.FinishOutpatientRequest{}

	if err := c.Bind(&outpatient); err != nil {
		errMsg = "Unable to parse request body"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindBadRequest))
	}

	if err := p.validate.Struct(outpatient); err != nil {
		errMsg = "Invalid request. Make sure all fields are filled correctly"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindUnprocessable))
	}

	err := p.business.FinishOutpatient(outpatient.ToOutpatientCore(), userID, role)
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}

	return response.Success(c, code, message, nil)
}

func (p *SchedulePresentation) DeleteOutpatient(c echo.Context) error {
	const op errors.Op = "schedules.presentation.DeleteOutpatient"
	var errMsg errors.ErrClientMessage

	code := http.StatusOK
	message := "Successfully deleting outpatient"

	outpatientID, err := strconv.Atoi(c.Param("outpatientId"))
	if err != nil {
		errMsg = "Invalid outpatient id"
		return response.Error(c, errors.E(err, op, errMsg, errors.KindBadRequest))
	}

	err = p.business.RemoveOutpatientById(outpatientID)
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}

	return response.Success(c, code, message, nil)
}
