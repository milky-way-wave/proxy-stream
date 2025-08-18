package datetime

import (
	"fmt"
	"time"
)

func Datetime2duration(t *time.Time) string {
	if t == nil {
		return ""
	}

	now := time.Now()
	elapsed := now.Sub(*t)

	if elapsed < 0 {
		elapsed = -elapsed
	}

	years := int(elapsed.Hours() / 24 / 365.25)
	remainingHours := elapsed.Hours() - float64(years)*365.25*24

	months := int(remainingHours / 24 / 30.44)
	remainingHours -= float64(months) * 30.44 * 24

	days := int(remainingHours / 24)
	remainingHours -= float64(days) * 24

	hours := int(remainingHours)
	minutes := int(elapsed.Minutes()) - years*365*24*60 - months*30*24*60 - days*24*60 - hours*60
	seconds := int(elapsed.Seconds()) - years*365*24*3600 - months*30*24*3600 - days*24*3600 - hours*3600 - minutes*60

	var duration string

	if years > 0 {
		duration += fmt.Sprintf("%d years ", years)
	}
	if months > 0 {
		duration += fmt.Sprintf("%d months ", months)
	}
	if days > 0 {
		duration += fmt.Sprintf("%dd ", days)
	}
	if hours > 0 {
		duration += fmt.Sprintf("%dh ", hours)
	}
	if minutes > 0 {
		duration += fmt.Sprintf("%dm ", minutes)
	}
	if seconds > 0 || duration == "" {
		duration += fmt.Sprintf("%ds", seconds)
	}

	return duration
}
