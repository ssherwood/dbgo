package randoms

import "math/rand"

//
// RandomBoolean returns a random true (or false) value about 50% of the time.
//
func RandomBoolean() bool {
	return RandomBooleanRatioOfFalse(1)
}

//
// RandomBooleanRatioOfFalse returns a boolean result based on a ratio of false
// to true.  A larger weight value returns a higher percentage change of being
// false.
//
// Weight:   0   1    2    3    4    5    6   . 100
// % False:  0%  50%  66%  75%  80%  83%  ... . 1%
//
func RandomBooleanRatioOfFalse(weight int) bool {
	return rand.Intn(weight+1) == 0
}

//
// RandomBooleanRatioOfTrue returns a boolean result based on a ratio of true
// to false.  A larger weight value returns a higher percentage change of being
// true.
//
// Weight:   0   1    2    3    4    5    6   . 100
// % True:   0%  50%  66%  75%  80%  83%  ... . 1%
//
func RandomBooleanRatioOfTrue(weight int) bool {
	return rand.Intn(weight+1) > 0
}
