package ads

import "time"

type Date struct {
	Day   int
	Month time.Month
	Year  int
}

func CurrentDate() Date {
	t := time.Now().UTC()
	return Date{
		Day:   t.Day(),
		Month: t.Month(),
		Year:  t.Year(),
	}
}
