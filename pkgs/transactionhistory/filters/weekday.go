package filters

import (
	"time"

	"4d63.com/clippercard/pkgs/transactionhistory"
)

// Weekday filters transactions by the given weekday, returning only
// transactions that occurred on that weekday.
func Weekday(transactions []transactionhistory.Transaction, weekdays []time.Weekday) []transactionhistory.Transaction {
	filtered := []transactionhistory.Transaction{}
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
