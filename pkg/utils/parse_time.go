package utils

import "time"

// ParseTime parses a string in the format "HH:MM" into a time.Time object.
func ParseTime(timeStr string) (*time.Time, error) {
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return nil, err
	}

	normalized := time.Date(
		1970, 1, 1,
		t.Hour(), t.Minute(), 0,
		0,
		time.UTC,
	)

	return &normalized, nil
}

// ParseDate parses a string in the format "YYYY-MM-DD" into a time.Time object.
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// pointer to time
func PtrTime(t time.Time) *time.Time {
	return &t
}
