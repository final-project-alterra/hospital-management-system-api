package schedules

type ScheduleQuery struct {
	StartDate string
	EndDate   string
	Limit     int
	Repeat    string
}

const (
	RepeatNoRepeat = "no-repeat"
	RepeatDaily    = "daily"
	RepeatWeekly   = "weekly"
	RepeatMonthly  = "monthly"

	StatusFinished   = 1
	StatusOnprogress = 2
	StatusWaiting    = 3
	StatusCanceled   = 4
)
