package request

import (
	"time"

	"github.com/final-project-alterra/hospital-management-system-api/features/schedules"
	"github.com/go-playground/validator/v10"
)

type QueryParamsRequest struct {
	StartDate string `query:"startdate" validate:"ValidateQueryDate"`
	EndDate   string `query:"enddate"`
	Limit     int    `query:"limit" validate:"ValidateQueryLimit"`
}

func NewQueryParamsRequest() QueryParamsRequest {
	return QueryParamsRequest{
		StartDate: "1900-01-01",
		EndDate:   "3000-01-01",
		Limit:     1000000,
	}
}

func (q QueryParamsRequest) ToScheduleQuery() schedules.ScheduleQuery {
	return schedules.ScheduleQuery{
		StartDate: q.StartDate,
		EndDate:   q.EndDate,
		Limit:     q.Limit,
	}
}

func ValidateQueryDate(fl validator.FieldLevel) bool {
	query, ok := fl.Parent().Interface().(QueryParamsRequest)
	if !ok {
		return false
	}

	start, err := time.Parse("2006-01-02", query.StartDate)
	if err != nil {
		return false
	}

	end, err := time.Parse("2006-01-02", query.EndDate)
	if err != nil {
		return false
	}

	return start.Before(end) || start.Equal(end)
}

func ValidateQueryLimit(fl validator.FieldLevel) bool {
	query, ok := fl.Parent().Interface().(QueryParamsRequest)
	if !ok {
		return false
	}

	if query.Limit < 0 {
		return false
	}

	return true
}
