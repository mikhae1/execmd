package execmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClusterSSHCmd(t *testing.T) {
	assert := assert.New(t)

	sshHosts := []string{"localhost", "127.0.0.1"}

	cluster := NewClusterSSHCmd(sshHosts)

	// parallel
	res, err := cluster.Run("VAR=world; echo Hello stdout $VAR; echo Hello stderr $VAR >&2")
	assert.NoError(err)
	assert.EqualValues(len(res), len(sshHosts), "number of results not equal to hosts number")
	for i := range res {
		assert.EqualValues(res[i].res.stdout.String(), "Hello stdout world"+LineBreak)
		assert.EqualValues(res[i].res.stderr.String(), "Hello stderr world"+LineBreak)
	}

	res, err = cluster.Run("give-me-error")
	assert.Error(err)
	assert.EqualValues(len(res), len(sshHosts), "number of results not equal to hosts number")
	for i := range res {
		assert.Contains(res[i].res.stderr.String(), "give-me-error")
	}

	// serial
	res, err = cluster.RunSeq("VAR=world; echo Hello stdout $VAR; echo Hello stderr $VAR >&2")
	assert.NoError(err)
	assert.EqualValues(len(res), len(sshHosts), "number of results not equal to hosts number")

	for i := range res {
		assert.EqualValues(res[i].res.stdout.String(), "Hello stdout world"+LineBreak)
		assert.EqualValues(res[i].res.stderr.String(), "Hello stderr world"+LineBreak)
	}

	cluster.StopOnError = true
	res, err = cluster.RunSeq("give-me-error")
	assert.Error(err)
	assert.EqualValues(len(res), 1, "more than one result")
	assert.Contains(res[0].res.stderr.String(), "give-me-error")

	h := []string{}
	h = append(h, sshHosts[0], sshHosts[0])
	twinCluster := NewClusterSSHCmd(h)
	res, err = twinCluster.RunSeq("FILE=.dp_remove_me; [[ -f $FILE ]] || (touch $FILE; echo 1) && (echo 2; rm $FILE)")
	assert.NoError(err)
	for i := range res {
		assert.EqualValues(res[i].res.stdout.String(), "1"+LineBreak+"2"+LineBreak)
	}
}
