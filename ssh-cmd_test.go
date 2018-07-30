package execmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSSHCmd(t *testing.T) {
	assert := assert.New(t)

	const sshHost = "localhost"

	srv := NewSSHCmd(sshHost)
	res, err := srv.Run("VAR=world; echo Hello stdout $VAR; echo Hello stderr $VAR >&2")
	assert.NoError(err)
	assert.EqualValues(res.stdout.String(), "Hello stdout world"+LineBreak)
	assert.EqualValues(res.stderr.String(), "Hello stderr world"+LineBreak)

	res, err = srv.Run("i-am-not-exist")
	assert.Error(err)
	assert.Contains(res.stderr.String(), "i-am-not-exist")
}
