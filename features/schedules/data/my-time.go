package data

import (
	"database/sql/driver"
	"strconv"
	"strings"
	"time"

	"github.com/final-project-alterra/hospital-management-system-api/errors"
)

const MyTimeFormat = "15:04:05"

type MyTime time.Time

/* Time format for gorm */

func NewMyTime(timeInput string) (MyTime, error) {
	const op errors.Op = "schedules.data.NewMyTime"
	var errMsg errors.ErrClientMessage = "Invalid time format"
	var t time.Time

	times := strings.Split(timeInput, ":")
	if len(times) != 3 {
		return MyTime(t), errors.E(errors.New(string(errMsg)), op, errMsg, errors.KindBadRequest)
	}

	hour, err := strconv.Atoi(times[0])
	if err != nil {
		return MyTime(t), errors.E(errors.New(string(errMsg)), op, errMsg, errors.KindBadRequest)
	}

	minute, err := strconv.Atoi(times[1])
	if err != nil {
		return MyTime(t), errors.E(errors.New(string(errMsg)), op, errMsg, errors.KindBadRequest)
	}

	second, err := strconv.Atoi(times[2])
	if err != nil {
		return MyTime(t), errors.E(errors.New(string(errMsg)), op, errMsg, errors.KindBadRequest)
	}

	t = time.Date(0, time.January, 1, hour, minute, second, 0, time.UTC)
	return MyTime(t), nil
}

func (t *MyTime) Scan(value interface{}) error {
	const op errors.Op = "schedules.data.(MyTime)Scan"
	var errMsg errors.ErrClientMessage = "Unsupported MyTime scan value type"

	switch v := value.(type) {
	case []byte:
		return t.ParseTime(string(v))
	case string:
		return t.ParseTime(v)
	case time.Time:
		*t = MyTime(v)
	case nil:
		*t = MyTime{}
	default:
		return errors.E(errors.New(string(errMsg)), op, errMsg, errors.KindBadRequest)
	}
	return nil
}

func (t MyTime) Value() (driver.Value, error) {
	return driver.Value(time.Time(t).Format(MyTimeFormat)), nil
}

func (t *MyTime) ParseTime(value string) error {
	dd, err := time.Parse(MyTimeFormat, value)
	if err != nil {
		return err
	}
	*t = MyTime(dd)
	return nil
}

func (t *MyTime) String() string {
	if time.Time(*t).IsZero() {
		return ""
	}
	return time.Time(*t).Format(MyTimeFormat)
}

func (MyTime) GormDataType() string {
	return "TIME"
}
