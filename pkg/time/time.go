package time

import (
	"time"
)

// getCurrentTime retrieves the current time and returns it as a formatted string.
// The time is formatted as "Monday, January 2, 2006 at 3:04pm".
func getCurrentTime() string {
	currentTime := time.Now()

	return currentTime.Format("Monday, January 2, 2006 at 3:04pm")
}
