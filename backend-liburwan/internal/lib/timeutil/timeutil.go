package timeutil

import (
	"time"
)

var Loc *time.Location

func init() {
	var err error
	Loc, err = time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Fallback to UTC+7 if location data is not available
		Loc = time.FixedZone("WIB", 7*3600)
	}
}

func Now() time.Time {
	return time.Now().In(Loc)
}

func StartOfToday() time.Time {
	now := Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, Loc)
}

func StartOfCurrentMonth() time.Time {
	now := Now()
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, Loc)
}

func ParseYYYYMM(bulan string) (time.Time, error) {
	t, err := time.ParseInLocation("2006-01", bulan, Loc)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
