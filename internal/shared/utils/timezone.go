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
