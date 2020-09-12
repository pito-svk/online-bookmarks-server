package entity

import "time"

func calcDurationInMs(t1 time.Time, t2 time.Time) int {
	return int(t1.Sub(t2).Milliseconds())
}
