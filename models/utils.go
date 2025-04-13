package models

import "time"

func sameDay(d1, d2 time.Time) bool {
	y1, m1, d1day := d1.Date()
	y2, m2, d2day := d2.Date()
	return y1 == y2 && m1 == m2 && d1day == d2day
}
