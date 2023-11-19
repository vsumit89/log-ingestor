package helper

import (
	"fmt"
	"time"
)

func GetFiltersUsingTime(startTime, endTime *time.Time) string {
	if startTime == nil && endTime == nil {
		return "" // Both are nil, no valid range
	}

	if startTime == nil {
		// Only endTime is valid
		endUnix := endTime.Unix()
		return fmt.Sprintf("timestamp <= %d", endUnix)
	}

	if endTime == nil {
		// Only startTime is valid
		startUnix := startTime.Unix()
		return fmt.Sprintf("timestamp >= %d", startUnix)
	}

	// Both startTime and endTime are valid
	startUnix := startTime.Unix()
	endUnix := endTime.Unix()

	// Format the filter string
	dateFilter := fmt.Sprintf("timestamp >= %d AND timestamp <= %d", startUnix, endUnix)

	return dateFilter
}
