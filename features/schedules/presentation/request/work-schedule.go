package request

import (
	"time"

	"github.com/final-project-alterra/hospital-management-system-api/features/schedules"
	"github.com/go-playground/validator/v10"
)

type CreateWorkScheduleRequest struct {
	DoctorID  int    `json:"doctorId" validate:"gt=0"`
	NurseID   int    `json:"nurseId" validate:"gt=0"`
	StartTime string `json:"startTime" validate:"required,ValidateCreateScheduleTime"`
	EndTime   string `json:"endTime" validate:"required"`

	StartDate string `json:"startDate" validate:"ValidateCreateScheduleDate"`
	EndDate   string `json:"endDate"`
	Repeat    string `json:"repeat" validate:"oneof='no-repeat' 'daily' 'weekly' 'monthly'"`
}

type UpdateWorkScheduleRequest struct {
	ID        int    `json:"id" validate:"gt=0"`
	DoctorID  int    `json:"doctorId" validate:"gt=0"`
	NurseID   int    `json:"nurseId" validate:"gt=0"`
	Date      string `json:"date" validate:"ValidateUpdateScheduleDate"`
	StartTime string `json:"startTime" validate:"required,ValidateUpdateScheduleTime"`
	EndTime   string `json:"endTime" validate:"required"`
}

func (w UpdateWorkScheduleRequest) ToWorkScheduleCore() schedules.WorkScheduleCore {
	wc := schedules.WorkScheduleCore{}
	wc.ID = w.ID
	wc.Doctor.ID = w.DoctorID
	wc.Nurse.ID = w.NurseID
	wc.Date = w.Date
	wc.StartTime = w.StartTime
	wc.EndTime = w.EndTime

	return wc
}

func ValidateCreateScheduleDate(fl validator.FieldLevel) bool {
	input, ok := fl.Parent().Interface().(CreateWorkScheduleRequest)
	if !ok {
		return false
	}

	if input.Repeat != schedules.RepeatNoRepeat && input.EndDate == "" {
		return false
	}

	start, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return false
	}

	end, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		return false
	}

	return start.Before(end) || start.Equal(end)
}

func ValidateCreateScheduleTime(fl validator.FieldLevel) bool {
	input, ok := fl.Parent().Interface().(CreateWorkScheduleRequest)
	if !ok {
		return false
	}

	startTime, err := time.Parse("15:04:05", input.StartTime)
	if err != nil {
		return false
	}

	endTime, err := time.Parse("15:04:05", input.EndTime)
	if err != nil {
		return false
	}

	return startTime.Before(endTime)
}

func ValidateUpdateScheduleTime(fl validator.FieldLevel) bool {
	input, ok := fl.Parent().Interface().(UpdateWorkScheduleRequest)
	if !ok {
		return false
	}

	startTime, err := time.Parse("15:04:05", input.StartTime)
	if err != nil {
		return false
	}

	endTime, err := time.Parse("15:04:05", input.EndTime)
	if err != nil {
		return false
	}

	return startTime.Before(endTime)
}
func ValidateUpdateScheduleDate(fl validator.FieldLevel) bool {
	input, ok := fl.Parent().Interface().(UpdateWorkScheduleRequest)
	if !ok {
		return false
	}

	_, err := time.Parse("2006-01-02", input.Date)
	return err == nil
}
