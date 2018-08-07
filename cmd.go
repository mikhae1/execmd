/*
	Wrapper for https://golang.org/pkg/os/exec/ to invoke command in shell,
	pipe stdout, stderr to console with prefixes,
	record output buffers

	Copyright(c) 2018 mink0
*/

package execmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// List of system shell binaries
var shellPathList = []string{os.Getenv("SHELL"), "bash", "sh"}

// Cmd is wrapper struct around exec.Cmd
// Stream buffers will be saved with Record[Stdout|Stderr] == true
// Mute options turns off console output
type Cmd struct {
	ShellPath    string
	RecordStdout bool
	RecordStderr bool
	MuteStdout   bool
	MuteStderr   bool
	MuteCmd      bool
	PrefixStdout string // stdout prining prefix
	PrefixStderr string // stderr priniting prefix
	PrefixCmd    string // command printing prefix

	Cmd *exec.Cmd // os.Exec instance
}

// CmdRes resulting struct
type CmdRes struct {
	Stdout *bytes.Buffer
	Stderr *bytes.Buffer
}

// NewCmd initializes Cmd with defaults
func NewCmd() *Cmd {
	return &Cmd{
		RecordStdout: true,
		RecordStderr: true,
		PrefixCmd:    "$ ",
		PrefixStdout: colorOK("> "),
		PrefixStderr: colorErr("@err "),
	}
}

// Wait is a exec.Wait wrapper with buffer flushes
func (c *Cmd) Wait() (err error) {
	err = c.Cmd.Wait()

	// flush last line
	c.Cmd.Stderr.(*pStream).Close()
	c.Cmd.Stdout.(*pStream).Close()
	return
}

// Run is exec.Run() wrapper: runs command and blocking wait for result
func (c *Cmd) Run(command string) (res CmdRes, err error) {
	if res, err = c.Start(command); err != nil {
		return
	}

	err = c.Wait()
	return
}

// Start is exec.Start() wrapper with system shell and output buffers initialization
func (c *Cmd) Start(command string) (res CmdRes, err error) {
	if c.ShellPath == "" {
		if c.ShellPath, err = findPath(shellPathList); err != nil {
			err = errors.Wrapf(err, "can't find shell binary: %v", shellPathList)
			return
		}
	}

	c.Cmd = exec.Command(c.ShellPath, "-c", command)

	// FIXME: rewrite to use raw buffers only when mute == true
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

	c.Cmd.Stdin = os.Stdin

	if !c.MuteCmd {
		fmt.Printf("%s%s\n", c.PrefixCmd, command)
	}

	err = c.Cmd.Start()

	res.Stdout = stdoutStream.Get()
	res.Stderr = stderrStream.Get()
	return
}

func findPath(paths []string) (path string, err error) {
	for _, p := range paths {
		path, err = exec.LookPath(p)
		if err == nil {
			break
		}
	}

	return
}
