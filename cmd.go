// Package execmd provides a wrapper around https://golang.org/pkg/os/exec/
// to execute commands in a shell, pipe stdout and stderr to the console with
// prefixes, and record output buffers.
package execmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
)

// Default shell paths that can be overridden by setting the `SHELL` environment variable.
var shellPathList = []string{os.Getenv("SHELL"), "bash", "sh"}

// Cmd is a wrapper struct around exec.Cmd that provides additional
// functionality such as recording and muting stdout and stderr, and
// customizing output prefixes.
type Cmd struct {
	ShellPath    string
	Interactive  bool
	LoginShell   bool
	RecordStdout bool
	RecordStderr bool
	MuteStdout   bool
	MuteStderr   bool
	MuteCmd      bool
	PrefixStdout string
	PrefixStderr string
	PrefixCmd    string

	Cmd *exec.Cmd
}

// CmdRes represents the result of a command, including the stdout and stderr buffers.
type CmdRes struct {
	Stdout *bytes.Buffer
	Stderr *bytes.Buffer
}

// NewCmd initializes a Cmd with default settings.
func NewCmd() *Cmd {
	cmd := Cmd{
		RecordStdout: true,
		RecordStderr: true,
		PrefixCmd:    "$ ",
		PrefixStdout: colorOK("> "),
		PrefixStderr: colorErr("@err "),
	}

	if shellPath, err := findPath(shellPathList); err == nil {
		cmd.ShellPath = shellPath
	}

	return &cmd
}

// Wait wraps exec.Wait and ensures that the buffers are flushed after waiting.
func (c *Cmd) Wait() error {
	err := c.Cmd.Wait()

	c.Cmd.Stderr.(*pStream).Close()
	c.Cmd.Stdout.(*pStream).Close()
	return err
}

// Run executes the command and waits for the result.
func (c *Cmd) Run(command string) (CmdRes, error) {
	res, err := c.Start(command)
	if err != nil {
		return res, err
	}

	err = c.Wait()
	return res, err
}

// Start initializes the system shell and output buffers, and starts the command.
func (c *Cmd) Start(command string) (CmdRes, error) {
	args := []string{}
	if c.Interactive {
		args = append(args, "-i")
	}

	if c.LoginShell {
		args = append(args, "-l")
	}

	args = append(args, "-c", command)

	c.Cmd = exec.Command(c.ShellPath, args...)

	stdoutLogFile := log.New(os.Stdout, "", 0)
	if c.MuteStdout {
		stdoutLogFile = log.New(bytes.NewBuffer([]byte("")), "", 0)
	}

	stderrLogFile := log.New(os.Stderr, "", 0)
	if c.MuteStderr {
		stderrLogFile = log.New(bytes.NewBuffer([]byte("")), "", 0)
	}

	stdoutStream := newPStream(stdoutLogFile, c.PrefixStdout, c.RecordStdout)
	c.Cmd.Stdout = stdoutStream

	stderrStream := newPStream(stderrLogFile, c.PrefixStderr, c.RecordStderr)
	c.Cmd.Stderr = stderrStream
	if c.Interactive {
		c.Cmd.Stdin = os.Stdin
	}

	if !c.MuteCmd {
		fmt.Printf("%s%s\n", c.PrefixCmd, colorStrong(command))
	}

	err := c.Cmd.Start()

	res := CmdRes{
		Stdout: stdoutStream.Get(),
		Stderr: stderrStream.Get(),
	}
	return res, err
}

// findPath searches for the available shell path from a given list of paths.
func findPath(paths []string) (string, error) {
	for _, p := range paths {
		path, err := exec.LookPath(p)
		if err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("no valid shell found in path list")
}
