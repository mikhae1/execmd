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
	assert.EqualValues(len(sshHosts), len(res), "number of results not equal to hosts number")
	for i := range res {
		assert.EqualValues("Hello stdout world\n", res[i].Res.Stdout.String())
		assert.EqualValues("Hello stderr world\n", res[i].Res.Stderr.String())
	}

	res, err = cluster.Run("give-me-error")
	assert.Error(err)
	assert.EqualValues(len(sshHosts), len(res), "number of results not equal to hosts number")
	for i := range res {
		assert.Contains(res[i].Res.Stderr.String(), "give-me-error")
	}

	// serial
	res, err = cluster.RunOneByOne("VAR=world; echo Hello stdout $VAR; echo Hello stderr $VAR >&2")
	assert.NoError(err)
	assert.EqualValues(len(sshHosts), len(res), "number of results not equal to hosts number")

	for i := range res {
		assert.EqualValues("Hello stdout world\n", res[i].Res.Stdout.String())
		assert.EqualValues("Hello stderr world\n", res[i].Res.Stderr.String())
	}

	cluster.StopOnError = true
	res, err = cluster.RunOneByOne("give-me-error")
	assert.Error(err)
	assert.EqualValues(1, len(res), "more than one result")
	assert.Contains(res[0].Res.Stderr.String(), "give-me-error")

	h := []string{}
	h = append(h, sshHosts[0], sshHosts[0])
	twinCluster := NewClusterSSHCmd(h)
	res, err = twinCluster.RunOneByOne("FILE=.dp_remove_me; [[ -f $FILE ]] || (touch $FILE; echo 1) && (echo 2; rm $FILE)")
	assert.NoError(err)
	for i := range res {
		assert.EqualValues("1\n2\n", res[i].Res.Stdout.String())
	}
}
