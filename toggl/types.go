package toggl

import "time"

// Int returns a pointer to the int value.
func Int(v int) *int {
	return &v
}

// Float64 returns a pointer to the float64 value.
func Float64(v float64) *float64 {
	return &v
}

// String returns a pointer to the string value.
func String(v string) *string {
	return &v
}

// Bool returns a pointer to the bool value.
func Bool(v bool) *bool {
	return &v
}

// Time returns a pointer to the time.Time value.
func Time(v time.Time) *time.Time {
	return &v
}
