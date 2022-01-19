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

	StatusOnprogress = 1
	StatusWaiting    = 2
	StatusFinished   = 3
	StatusCanceled   = 4
)
