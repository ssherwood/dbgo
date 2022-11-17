package randoms

import "math/rand"

func RandomIntBetween(maxValue int, minValue int) int {
	return rand.Intn(maxValue-minValue+1) + minValue
}
