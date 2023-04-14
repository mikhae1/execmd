package execmd_test

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"

	execmd "github.com/mink0/exec-cmd.git"
)

func TestInteractive(t *testing.T) {
	cmd := execmd.NewCmd()
	cmd.Interactive = true

	res, err := cmd.Run("echo Hello stdout $USER; echo Hello stderr $USER >&2")
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}

	expected := "Hello stdout " + os.Getenv("USER") + "\n"
	if res.Stdout.String() != expected {
		t.Errorf("Unexpected output: %s", res.Stdout.String())
	}
}

func TestNewCmd(t *testing.T) {
	cmd := execmd.NewCmd()
	res, err := cmd.Run("echo Hello stdout $USER; echo Hello stderr $USER >&2")
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}

	expected := "Hello stdout " + os.Getenv("USER") + "\n"
	if res.Stdout.String() != expected {
		t.Errorf("Unexpected stdout output: %s", res.Stdout.String())
	}

	expected = "Hello stderr " + os.Getenv("USER") + "\n"
	if res.Stderr.String() != expected {
		t.Errorf("Unexpected stderr output: %s", res.Stderr.String())
	}

	cmd = execmd.NewCmd()
	res, err = cmd.Run("i-am-not-exist")
	if err == nil {
		t.Fatalf("Expected error but got none")
	}

	if !strings.Contains(res.Stderr.String(), "i-am-not-exist") {
		t.Errorf("Unexpected error output: %s", res.Stderr.String())
	}
}

func TestRunWithTimeout(t *testing.T) {
	cmd := execmd.NewCmd()

	// Test running a command with a timeout
	res, err := cmd.Run("sleep 1; echo OK", 3*time.Second)
	if err != nil {
		t.Fatalf("Unexpected error due to timeout: %v, %s", res, err)
	}
	if !strings.Contains(res.Stdout.String(), "OK") {
		t.Errorf("Unexpected output: %s", res.Stdout.String())
	}

	// Test running a command that will be killed due to timeout
	res, err = cmd.Run("sleep 3; echo OK", 1*time.Second)
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				if status.Signal() != syscall.SIGKILL {
					t.Errorf("Unexpected error output: %s", err)
				}
			}
		} else {
			t.Errorf("Unexpected error output: %s", err)
		}
	} else {
		t.Errorf("Expected an error due to the process being killed by a timeout")
	}
}
