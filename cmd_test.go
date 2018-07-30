package execmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCmd(t *testing.T) {
	assert := assert.New(t)

	cmd := NewCmd()
	res, err := cmd.Run("echo Hello stdout $USER; echo Hello stderr $USER >&2")
	assert.NoError(err)
	assert.EqualValues(res.Stdout.String(), "Hello stdout "+os.Getenv("USER")+"\n")
	assert.EqualValues(res.Stderr.String(), "Hello stderr "+os.Getenv("USER")+"\n")

	cmd = NewCmd()
	res, err = cmd.Run("i-am-not-exist")
	assert.Error(err)
	assert.Contains(res.Stderr.String(), "i-am-not-exist")
}
