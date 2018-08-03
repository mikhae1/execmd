/*
	Wrapper for Cmd to invoke ssh commands via OpenSSH binary

	Copyright(c) 2018 mink0
*/

// TODO:
// User detect from user@host

package execmd

import (
	"os"
	"strings"

	"github.com/pkg/errors"
)

// SSHCmd is a wrapper on Cmd
type SSHCmd struct {
	Cmd *Cmd

	SSHBinPath  string
	Host        string
	User        string
	Interactive bool
	Port        string
	Key         string
	Cwd         string
}

// Path to ssh binary
// Could be overriden by setting SSH_BIN_PATH env variable
var sshBinList = []string{os.Getenv("SSH_BIN_PATH"), "ssh"}

// NewSSHCmd initialize SSHCmd with defaults
func NewSSHCmd(host string) *SSHCmd {
	ssh := &SSHCmd{
		Host: host,
	}

	ssh.Cmd = NewCmd()
	ssh.Cmd.PrefixStdout = color(host) + " "
	ssh.Cmd.PrefixStderr = color(host) + colorErr("@err ")
	ssh.Cmd.PrefixCmd = "$ "
	return ssh
}

// Wait wraps Cmd.Wait()
func (s *SSHCmd) Wait() error {
	return s.Cmd.Wait()
}

// Run wraps Cmd.Run()
func (s *SSHCmd) Run(command string) (res CmdRes, err error) {
	return s.Cmd.Run(command)
}

// Start wraps Cmd.Start() with ssh invocation
func (s *SSHCmd) Start(command string) (res CmdRes, err error) {
	if s.SSHBinPath == "" {
		if s.SSHBinPath, err = findPath(sshBinList); err != nil {
			err = errors.Wrapf(err, "can't find ssh binary: %v", shellPathList)
			return
		}
	}

	sshArgs := s.warpInSSH(command)

	res, err = s.Cmd.Start(strings.Join(sshArgs, " "))

	return
}

// transform `command` into ssh-compatible argument string
func (s *SSHCmd) warpInSSH(command string) (sshArgs []string) {
	sshArgs = append(sshArgs, s.SSHBinPath)

	hostWithUser := s.Host
	if s.User != "" {
		hostWithUser = s.User + "@" + s.Host
	}

	sshArgs = append(sshArgs, hostWithUser)

	if strings.Contains(command, "sudo") || s.Interactive {
		sshArgs = append(sshArgs, "-tt")
	}
	if s.Port != "" {
		sshArgs = append(sshArgs, "-p", s.Port)
	}
	if s.Key != "" {
		sshArgs = append(sshArgs, "-i", s.Key)
	}
	if s.Cwd != "" {
		command = "cd " + s.Cwd + " && " + command
	}

	// escape single quotes for shell encapsulation
	sshArgs = append(sshArgs, "'"+strings.Replace(command, "'", "'\\''", -1)+"'")

	return
}
