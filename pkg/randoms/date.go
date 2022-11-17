package randoms

import (
	"math/rand"
	"time"
)

func RandomDateUTC(maxYear int, minYear int) time.Time {
	if minYear > maxYear {
		minYear = maxYear - 1
	}

	min := time.Date(minYear, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(maxYear, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	return time.Unix(rand.Int63n(max-min)+min, 0).UTC()
}
