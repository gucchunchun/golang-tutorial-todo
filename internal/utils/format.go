package utils

import (
	"fmt"
	"time"
)

const (
	DateFormatOutput     = "2006-01-02"
	DatetimeFormatOutput = DateFormatOutput + " 15:04:05"
)

func FormatTaskOutput(id uint, name, status string, createdAt, dueDate time.Time) string {
	createdAtOutput, _ := formatDatetime(createdAt)
	dueDateOutput, _ := formatDate(dueDate)
	timeLeft, _ := formatTimeLeft(dueDate)

	// 左揃え＋固定幅でフォーマット
	return fmt.Sprintf("%-4d | %-20s | %-8s | %-19s | %-10s | %-15s",
		id, name, status, createdAtOutput, dueDateOutput, timeLeft)
}

func formatDate(date time.Time) (string, error) {
	if date.IsZero() {
		return "No due date", nil
	}
	return date.Format(DateFormatOutput), nil
}

func formatDatetime(datetime time.Time) (string, error) {
	if datetime.IsZero() {
		return "No date", nil
	}
	return datetime.Format(DatetimeFormatOutput), nil
}

func formatTimeLeft(dueDate time.Time) (string, error) {
	if dueDate.IsZero() {
		return "No due date", nil
	}

	now := time.Now()
	if dueDate.Before(now) {
		return "Overdue", nil
	}

	timeLeft := dueDate.Sub(now)
	days := int(timeLeft.Hours() / 24)
	hours := int(timeLeft.Hours()) % 24

	return fmt.Sprintf("%d days %d hours", days, hours), nil
}
