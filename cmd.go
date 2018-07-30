/*
	Wrapper for https://golang.org/pkg/os/exec/ to invoke command in shell,
	pipe stdout, stderr to console with prefixes,
	record output buffers

	Copyright(c) 2018 mink0
*/

// TODO: add exit code capture

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

// Cmd wrapper struct around exec.Cmd
// Stream buffers will be saved with Record[Stdout|Stderr] == true
// Mute option turns off output
type Cmd struct {
	cmd          *exec.Cmd
	ShellPath    string
	Prefix       CmdPrefix
	RecordStdout bool
	RecordStderr bool
	MuteStdout   bool
	MuteCmd      bool
}

// CmdPrefix prefixes for stdout and stderr output.
// Cmd prefix is for command output
type CmdPrefix struct {
	cmd    string
	stdout string
	stderr string
}

// CmdRes resulting struct
type CmdRes struct {
	stdout *bytes.Buffer
	stderr *bytes.Buffer
	cmd    *exec.Cmd
}

// NewCmd Cmd constructor
func NewCmd() *Cmd {
	return &Cmd{
		RecordStdout: true,
		RecordStderr: true,
		Prefix: CmdPrefix{
			cmd:    "$ ",
			stdout: "> ",
			stderr: ColorErr("@err "),
		},
	}
}

func (c *Cmd) Wait() (err error) {
	err = c.cmd.Wait()

	// flush last line
	c.cmd.Stderr.(*PStream).Close()
	c.cmd.Stdout.(*PStream).Close()
	return
}

func (c *Cmd) Run(command string) (res CmdRes, err error) {
	if res, err = c.Start(command); err != nil {
		return
	}

	err = c.Wait()
	return
}

func FindPath(paths []string) (path string, err error) {
	for _, p := range paths {
		path, err = exec.LookPath(p)
		if err == nil {
			break
		}
	}

	return
}

func (c *Cmd) Start(command string) (res CmdRes, err error) {
	if c.ShellPath == "" {
		if c.ShellPath, err = FindPath(shellPathList); err != nil {
			err = errors.Wrapf(err, "can't find shell binary: %v", shellPathList)
			return
		}
	}

	c.cmd = exec.Command(c.ShellPath, "-c", command)

	// FIXME: rewrite to use raw buffers only when mute == true
	stdoutLogFile := log.New(os.Stdout, "", 0)
	if c.MuteStdout {
		stdoutLogFile = log.New(bytes.NewBuffer([]byte("")), "", 0)
	}

	stdoutStream := NewPStream(stdoutLogFile, c.Prefix.stdout, c.RecordStdout)
	c.cmd.Stdout = stdoutStream

	stderrStream := NewPStream(log.New(os.Stderr, "", 0), c.Prefix.stderr, c.RecordStderr)
	c.cmd.Stderr = stderrStream

	c.cmd.Stdin = os.Stdin

	if !c.MuteCmd {
		fmt.Printf("%s%s\n", c.Prefix.cmd, command)
	}

	err = c.cmd.Start()

	res.cmd = c.cmd
	res.stdout = stdoutStream.Get()
	res.stderr = stderrStream.Get()
	return
}
