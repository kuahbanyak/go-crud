package utils

import "time"

// WIBLocation is the UTC+7 timezone (Western Indonesian Time)
var WIBLocation *time.Location

func init() {
	var err error
	// Load UTC+7 timezone (Asia/Jakarta - Western Indonesian Time)
	WIBLocation, err = time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Fallback to fixed offset if timezone data not available
		WIBLocation = time.FixedZone("WIB", 7*60*60) // UTC+7
	}
}

// FormatTimeWIB formats a time.Time to RFC3339 string in WIB timezone (UTC+7)
func FormatTimeWIB(t time.Time) string {
	return t.In(WIBLocation).Format(time.RFC3339)
}

// NowWIB returns current time in WIB timezone (UTC+7)
func NowWIB() time.Time {
	return time.Now().In(WIBLocation)
}

// ParseTimeWIB parses a string to time.Time assuming WIB timezone
func ParseTimeWIB(value string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, err
	}
	return t.In(WIBLocation), nil
}

// GetWeekStart returns the start of the week (Monday 00:00:00) for the given date
func GetWeekStart(t time.Time) time.Time {
	// Get the weekday (0 = Sunday, 1 = Monday, ...)
	weekday := int(t.Weekday())

	// Calculate days to subtract to get to Monday
	// If Sunday (0), go back 6 days; if Monday (1), go back 0 days, etc.
	daysToMonday := (weekday + 6) % 7

	// Get Monday of the week
	monday := t.AddDate(0, 0, -daysToMonday)

	// Set to start of day (00:00:00)
	return time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, t.Location())
}

// GetWeekEnd returns the end of the week (Sunday 23:59:59) for the given date
func GetWeekEnd(t time.Time) time.Time {
	weekStart := GetWeekStart(t)
	// Add 6 days to Monday to get Sunday, then set to end of day
	sunday := weekStart.AddDate(0, 0, 6)
	return time.Date(sunday.Year(), sunday.Month(), sunday.Day(), 23, 59, 59, 999999999, t.Location())
}

// IsInCurrentOrFutureWeek checks if the given date is in the current week or a future week
func IsInCurrentOrFutureWeek(checkDate time.Time, referenceDate time.Time) bool {
	checkWeekStart := GetWeekStart(checkDate)
	referenceWeekStart := GetWeekStart(referenceDate)

	// Allow if checkDate's week is same or after reference week
	return !checkWeekStart.Before(referenceWeekStart)
}

// GetWeekNumber returns the ISO week number and year
func GetWeekNumber(t time.Time) (year, week int) {
	return t.ISOWeek()
}
