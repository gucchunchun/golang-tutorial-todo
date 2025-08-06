package utils

import "time"

const DateFormat = "2006-01-02"

func Now() time.Time {
	return time.Now()
}

func ParseDate(date string) (time.Time, error) {
	return time.Parse(DateFormat, date)
}
