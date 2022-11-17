package randoms

import (
	"math/rand"
	"strconv"
)

//
// RandomIPv4
// Returns a stupidly simple IPv4 address that is likely invalid or suspect, use only for generating fake data.
//
func RandomIPv4() string {
	return strconv.Itoa(rand.Intn(256)) + "." + strconv.Itoa(rand.Intn(256)) + "." + strconv.Itoa(rand.Intn(256)) + "." + strconv.Itoa(rand.Intn(256))
}
