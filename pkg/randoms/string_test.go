package randoms

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomStringAlpha(t *testing.T) {
	s := RandomStringAlpha(100)

	assert.Equal(t, len(s), 100)
}

func TestRandomStringAlphaNumeric(t *testing.T) {
	s := RandomStringAlphaNumeric(50)

	assert.Equal(t, len(s), 50)
}

func TestRandomStringAlphaNumericRandomLength(t *testing.T) {
	for i := 0; i < 100; i++ {
		s := RandomStringAlphaNumericRandomLength(50, 1)
		assert.True(t, len(s) >= 1 && len(s) < 50, "Length was %d", len(s))
	}
}

func TestRandomStringPassword(t *testing.T) {
	pwd := RandomStringPassword(16)

	assert.Equal(t, len(pwd), 16)
}
