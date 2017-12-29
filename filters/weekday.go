package filters

import (
	"time"

	"4d63.com/clippercardtransactionhistory"
)

// Weekday filters transactions by the given weekday, returning only
// transactions that occurred on that weekday.
func Weekday(transactions []clippercardtransactionhistory.Transaction, weekdays []time.Weekday) []clippercardtransactionhistory.Transaction {
	filtered := []clippercardtransactionhistory.Transaction{}
	for _, t := range transactions {
		if weekdaysContains(weekdays, t.Timestamp.Weekday()) {
			filtered = append(filtered, t)
		}
	}
	return filtered
}

func weekdaysContains(weekdays []time.Weekday, wd time.Weekday) bool {
	for _, weekday := range weekdays {
		if weekday == wd {
			return true
		}
	}
	return false
}
