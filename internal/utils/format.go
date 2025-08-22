package utils

import (
	"fmt"
	"time"

	"golang/tutorial/todo/internal/models"
)

const (
	DateFormatOutput     = "2006-01-02"
	DatetimeFormatOutput = DateFormatOutput + " 15:04:05"
)

func FormatDate(date models.Date) (string, error) {
	if date.IsZero() {
		return "", nil
	}
	return date.In(time.Local).Format(DateFormatOutput), nil
}

func FormatDatetime(datetime models.Date) (string, error) {
	if datetime.IsZero() {
		return "", nil
	}
	return datetime.In(time.Local).Format(DatetimeFormatOutput), nil
}

func FormatDurationToDays(duration time.Duration) string {
	if duration == 0 {
		return ""
	}

	days := int(duration.Hours() / 24)
	if int(duration.Hours())%24 != 0 {
		days++
	}

	return fmt.Sprintf("%d days", days)
}
