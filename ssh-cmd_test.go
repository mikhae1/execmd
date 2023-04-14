// You should enable `SSH` server locally and add your personal ssh key to `known_hosts` to avoid password prompting

package execmd_test

import (
	"strings"
	"testing"
	"time"

	execmd "github.com/mink0/exec-cmd.git"
)

const dummyHost = "localhost"

func TestNewSSHCmd_Run(t *testing.T) {
	srv := execmd.NewSSHCmd(dummyHost)
	res, err := srv.Run("VAR=world; echo Hello stdout $VAR; echo Hello stderr $VAR >&2")
	if err != nil {
		t.Error("Expected a timeout error, but got nil")
	}
	if res.Stdout.String() != "Hello stdout world\n" {
		t.Errorf("Expected and actual output do not match")
	}
	if res.Stderr.String() != "Hello stderr world\n" {
		t.Errorf("Expected and actual error output do not match")
	}
}

func TestNewSSHCmd_RunWithError(t *testing.T) {
	srv := execmd.NewSSHCmd(dummyHost)
	res, err := srv.Run("i-am-not-exist")
	if err == nil {
		t.Error("Expected error, but got nil")
	}
	if !strings.Contains(res.Stderr.String(), "i-am-not-exist") {
		t.Error("Expected error message not found")
	}
}

func TestNewSSHCmd_RunWithTimeout(t *testing.T) {
	srv := execmd.NewSSHCmd(dummyHost)
	res, err := srv.Run("sleep 3; echo OK", 1*time.Second)
	if err == nil {
		t.Error("Expected a timeout error, but got nil")
	}

	res, err = srv.Run("sleep 1; echo OK", 3*time.Second)
	if err != nil {
		t.Error("Unexpected a timeout error: %w", err)
	}
	if res.Stdout.String() != "OK\n" {
		t.Errorf("Expected and actual output do not match %s", res.Stdout)
	}
}

func TestNewSSHCmd_Cwd(t *testing.T) {
	srv := execmd.NewSSHCmd(dummyHost)
	srv.Cwd = "/tmp"
	res, err := srv.Run("pwd")
	if err != nil {
		t.Fatal(err)
	}
	if res.Stdout.String() != "/tmp\n" {
		t.Errorf("no working dir change")
	}
}

func TestNewSSHCmd_CwdNonExisting(t *testing.T) {
	srv := execmd.NewSSHCmd(dummyHost)
	srv.Cwd = "/i-am-nowhere"
	res, err := srv.Run("pwd")
	if err == nil {
		t.Error("Expected error, but got nil")
	}
	if !strings.Contains(res.Stderr.String(), "/i-am-nowhere") {
		t.Error("no error when nonexisting working dir change")
	}
}
