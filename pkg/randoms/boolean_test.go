package randoms

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

const (
	startingSeed = 1234567
)

//
// Note: these tests all rely on seeding rand with a known value (startingSeed)
// in order to be predictable and repeatable.  If adding additional assertions
// to existing tests, append them as to not change the RNG state.
//

func TestRandomBoolean(t *testing.T) {
	rand.Seed(startingSeed)

	b := RandomBoolean()
	assert.True(t, b, "RandomBoolean with startingSeed should always return true")
}

func TestRandomBooleanRatioOfFalse(t *testing.T) {

	// no seed change should be required, this should always return true
	b := RandomBooleanRatioOfFalse(0)
	assert.True(t, b, "RandomBooleanRatioOfFalse with weight of 0 should ALWAYS return true")

	rand.Seed(startingSeed)

	b2 := RandomBooleanRatioOfTrue(1)
	assert.False(t, b2, "RandomBooleanRatioOfTrue with a weight of 1 should be false with startingSeed")

	b3 := RandomBooleanRatioOfFalse(100)
	assert.False(t, b3, "RandomBooleanRationFalse with a weight of 100 should always return true with startingSeed")

	numberOfFalseResults := 0
	for i := 0; i < 1000; i++ {
		if !RandomBooleanRatioOfFalse(4) {
			numberOfFalseResults++
		}
	}
	assert.Equal(t, 793, numberOfFalseResults, "1k iterations of RandomBooleanRatioOfFalse with weight of 4 should be 793 (about 80%)")
}

func TestRandomBooleanRatioOfTrue(t *testing.T) {

	// no seed change should be required, this should always return false
	b := RandomBooleanRatioOfTrue(0)
	assert.False(t, b, "RandomBooleanRatioOfTrue with weight of 0 should ALWAYS return false")

	rand.Seed(startingSeed)

	b2 := RandomBooleanRatioOfTrue(1)
	assert.False(t, b2, "RandomBooleanRatioOfTrue with a weight of 1 should be false with the startingSeed")

	b3 := RandomBooleanRatioOfTrue(100)
	assert.True(t, b3, "RandomBooleanRatioOfTrue with a weight of 100 should be true with the startingSeed")

	numberOfTrueResults := 0
	for i := 0; i < 1000; i++ {
		if RandomBooleanRatioOfTrue(4) {
			numberOfTrueResults++
		}
	}
	assert.Equal(t, 793, numberOfTrueResults, "1k iterations of RandomBooleanRatioOfTrue with weight of 4 should be 793 (about 80%)")
}
