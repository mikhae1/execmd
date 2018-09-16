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
	res, err := cluster.Run("VAR=world; echo Parallel stdout $VAR; echo Parallel stderr $VAR >&2")
	assert.NoError(err)
	assert.EqualValues(len(sshHosts), len(res), "Run: number of results not equals to hosts number")
	for i := range res {
		assert.EqualValues("Parallel stdout world\n", res[i].Res.Stdout.String())
		assert.EqualValues("Parallel stderr world\n", res[i].Res.Stderr.String())
	}

	res, err = cluster.Run("give-me-error")
	assert.Error(err)
	assert.EqualValues(len(sshHosts), len(res), "Run with error: number of results not equals to hosts number")
	for i := range res {
		assert.Contains(res[i].Res.Stderr.String(), "give-me-error")
	}

	// serial
	res, err = cluster.RunOneByOne("VAR=world; echo Serial stdout $VAR; echo Serial stderr $VAR >&2")
	assert.NoError(err)
	assert.EqualValues(len(sshHosts), len(res), "RunOneByOne: number of results not equals to hosts number")

	for i := range res {
		assert.EqualValues("Serial stdout world\n", res[i].Res.Stdout.String())
		assert.EqualValues("Serial stderr world\n", res[i].Res.Stderr.String())
	}

	cluster.StopOnError = true
	res, err = cluster.RunOneByOne("give-me-error")
	assert.Error(err)
	assert.EqualValues(1, len(res), "RunOneByOne with stop on error: more than one result returned")
	assert.Contains(res[0].Res.Stderr.String(), "give-me-error")

	h := []string{}
	h = append(h, sshHosts[0], sshHosts[0])
	twinCluster := NewClusterSSHCmd(h)
	res, err = twinCluster.RunOneByOne("FILE=.dp_remove_me; [[ -f $FILE ]] || (touch $FILE; echo 1) && (echo 2; rm $FILE)")
	assert.NoError(err)
	for i := range res {
		assert.EqualValues("1\n2\n", res[i].Res.Stdout.String())
	}

	// check results are saved between executions
	res1, err1 := cluster.Run("echo res1")
	assert.NoError(err1)
	res2, err2 := cluster.Run("echo res2")
	assert.NoError(err2)
	assert.EqualValues(len(sshHosts), len(res), "number of results not equals to hosts number")
	for i := range res1 {
		assert.EqualValues("res1\n", res1[i].Res.Stdout.String())
		assert.EqualValues("res2\n", res2[i].Res.Stdout.String())
	}

	// check cwd changing
	cluster.Cwd = "/tmp"
	res, err = cluster.Run("pwd")
	assert.NoError(err)
	for i := range res {
		assert.EqualValues("/tmp\n", res[i].Res.Stdout.String(), "no working dir change")
	}
}
