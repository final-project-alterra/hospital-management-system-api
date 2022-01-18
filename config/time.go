package config

import (
	"errors"
	"time"
)

var timeLoc *time.Location

func InitTimeLoc(zone string) {
	if timeLoc != nil {
		return
	}

	t, err := time.LoadLocation(zone)
	if err != nil {
		panic(err.Error())
	}

	timeLoc = t
}

// It panics if `timeLoc` is nil
func GetTimeLoc() *time.Location {
	if timeLoc == nil {
		panic(errors.New("timeLoc is nil! make sure to call it in main"))
	}
	return timeLoc
}
