// You should enable `SSH` server locally and add your personal ssh key to `known_hosts` to avoid password prompting

package execmd_test

import (
	"strings"
	"testing"
	"time"

	execmd "github.com/mikhae1/execmd.git"
)

var dummyHosts = []string{"localhost", "127.0.0.1"}

func TestClusterSSHCmd_Run(t *testing.T) {
	cluster := execmd.NewClusterSSHCmd(dummyHosts)

	results, err := cluster.Run("VAR=world; echo Parallel stdout $VAR; echo Parallel stderr $VAR >&2")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(results) != len(dummyHosts) {
		t.Errorf("Run: number of results not equals to hosts number")
	}

	for _, res := range results {
		if res.Err != nil {
			t.Errorf("Error on host %s: %v", res.Host, res.Err)
		}
		if res.Res.Stdout.String() != "Parallel stdout world\n" {
			t.Errorf("Expected and actual output do not match for stdout")
		}
		if res.Res.Stderr.String() != "Parallel stderr world\n" {
			t.Errorf("Expected and actual error output do not match for stderr")
		}
	}
}

func TestClusterSSHCmd_StartAndWait(t *testing.T) {
	cluster := execmd.NewClusterSSHCmd(dummyHosts)
	results, err := cluster.Start("echo 'Hello, World!'")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	err = cluster.Wait()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	for _, res := range results {
		if res.Err != nil {
			t.Errorf("Error on host %s: %v", res.Host, res.Err)
		}

		if res.Res.Stdout.String() != "Hello, World!\n" {
			t.Errorf("Unexpected stdout on host %s: %s", res.Host, res.Res.Stdout)
		}
	}
}

func TestClusterSSHCmd_RunOneByOneWithTimeout(t *testing.T) {
	cluster := execmd.NewClusterSSHCmd(dummyHosts)

	cluster.StopOnError = false
	results, err := cluster.RunOneByOne("sleep 3; echo OK", 1*time.Second)
	if err != nil {
		t.Error("Unexpected timeout error: %w", err)
	}

	for _, res := range results {
		if res.Err == nil {
			t.Errorf("Expected an error on host %s, but got nil", res.Host)
		}
	}

	cluster.StopOnError = true
	_, err = cluster.RunOneByOne("sleep 3; echo OK", 1*time.Second)
	if err == nil {
		t.Error("Expected timeout error: %w", err)
	}

	cluster.StopOnError = true
	_, err = cluster.RunOneByOne("sleep 1; echo OK", 3*time.Second)
	if err != nil {
		t.Error("Unxpected timeout error: %w", err)
	}
}

func TestClusterSSHCmd_RunWithTimeout(t *testing.T) {
	cluster := execmd.NewClusterSSHCmd(dummyHosts)

	results, err := cluster.Run("sleep 3; echo OK", 1*time.Second)
	if err == nil {
		t.Error("Expected timeout error: %w", err)
	}

	for _, res := range results {
		if res.Err == nil {
			t.Errorf("Expected an error on host %s, but got nil", res.Host)
		}
	}

	_, err = cluster.Run("sleep 3; echo OK", 1*time.Second)
	if err == nil {
		t.Error("Expected timeout error: %w", err)
	}

	_, err = cluster.Run("sleep 1; echo OK", 3*time.Second)
	if err != nil {
		t.Error("Unxpected timeout error: %w", err)
	}
}

func TestNewClusterSSHCmd_RunWithError(t *testing.T) {
	cluster := execmd.NewClusterSSHCmd(dummyHosts)

	res, err := cluster.Run("give-me-error")
	if err == nil {
		t.Error("Expected error, but got nil")
	}
	if len(res) != len(dummyHosts) {
		t.Errorf("Run with error: number of results not equals to hosts number")
	}
	for i := range res {
		if !strings.Contains(res[i].Res.Stderr.String(), "give-me-error") {
			t.Error("Expected error message not found")
		}
	}
}

func TestNewClusterSSHCmd_RunOneByOne(t *testing.T) {
	cluster := execmd.NewClusterSSHCmd(dummyHosts)

	res, err := cluster.RunOneByOne("VAR=world; echo Serial stdout $VAR; echo Serial stderr $VAR >&2")
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != len(dummyHosts) {
		t.Errorf("RunOneByOne: number of results not equals to hosts number")
	}

	for i := range res {
		if res[i].Res.Stdout.String() != "Serial stdout world\n" {
			t.Errorf("Expected and actual output do not match")
		}
		if res[i].Res.Stderr.String() != "Serial stderr world\n" {
			t.Errorf("Expected and actual error output do not match")
		}
	}
}

func TestNewClusterSSHCmd_RunOneByOneStopOnError(t *testing.T) {
	cluster := execmd.NewClusterSSHCmd(dummyHosts)
	cluster.StopOnError = true

	res, err := cluster.RunOneByOne("give-me-error")
	if err == nil {
		t.Error("Expected error, but got nil")
	}
	if len(res) != 1 {
		t.Errorf("RunOneByOne with stop on error: more than one result returned")
	}
	if !strings.Contains(res[0].Res.Stderr.String(), "give-me-error") {
		t.Error("Expected error message not found")
	}
}

func TestNewClusterSSHCmd_CwdChange(t *testing.T) {
	sshHosts := dummyHosts
	cluster := execmd.NewClusterSSHCmd(sshHosts)
	cluster.Cwd = "/tmp"

	res, err := cluster.Run("pwd")
	if err != nil {
		t.Fatal(err)
	}
	for i := range res {
		if res[i].Res.Stdout.String() != "/tmp\n" {
			t.Errorf("no working dir change")
		}
	}
}

func TestNewClusterSSHCmd_MultiRun(t *testing.T) {
	sshHosts := dummyHosts
	cluster := execmd.NewClusterSSHCmd(sshHosts)

	res1, err1 := cluster.Run("echo res1")
	if err1 != nil {
		t.Fatal(err1)
	}
	res2, err2 := cluster.Run("echo res2")
	if err2 != nil {
		t.Fatal(err2)
	}
	if len(res1) != len(sshHosts) || len(res2) != len(sshHosts) {
		t.Errorf("number of results not equals to hosts number")
	}
	for i := range res1 {
		if res1[i].Res.Stdout.String() != "res1\n" {
			t.Errorf("Expected and actual output do not match")
		}
		if res2[i].Res.Stdout.String() != "res2\n" {
			t.Errorf("Expected and actual output do not match")
		}
	}
}
