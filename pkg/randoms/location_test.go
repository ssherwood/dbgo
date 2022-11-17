package randoms

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomMapEntry(t *testing.T) {
	var m = MapInterface{"KEY": {"VALUE"}}

	result := m.randomEntry("KEY", "DEFAULT")

	assert.Equal(t, "VALUE", result)
}
