package main

import (
	"fmt"
	"strings"
	"time"
)

const allWeekdays = "monday,tuesday,wednesday,thursday,friday,saturday,sunday"

func parseWeekdays(s string) ([]time.Weekday, error) {
	parts := strings.Split(s, ",")
	weekdays := make([]time.Weekday, 0, len(parts))
	for _, p := range parts {
		wd, err := parseWeekday(p)
		if err != nil {
			return nil, err
		}
		weekdays = append(weekdays, wd)
	}
	return weekdays, nil
}

func parseWeekday(s string) (time.Weekday, error) {
	weekdays := map[string]time.Weekday{
		"sunday":    time.Sunday,
		"monday":    time.Monday,
		"tuesday":   time.Tuesday,
		"wednesday": time.Wednesday,
		"thursday":  time.Thursday,
		"friday":    time.Friday,
		"saturday":  time.Saturday,
	}
	if wd, ok := weekdays[strings.ToLower(s)]; ok {
		return wd, nil
	}
	return time.Sunday, fmt.Errorf("weekday %q is not recognized", s)
}
