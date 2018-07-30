package execmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColor(t *testing.T) {
	str := "i am green"

	green := ColorOK(str)

	assert.Contains(t, green, str)
}
