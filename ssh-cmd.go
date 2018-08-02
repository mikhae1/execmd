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

type SSHCmd struct {
	Cmd *Cmd

	SshBinPath  string
	Host        string
	User        string
	Interactive bool
	Port        string
	Key         string
	Cwd         string
}

// Path to ssh binary
var sshBinList = []string{os.Getenv("SSH_BIN_PATH"), "ssh"}

func NewSSHCmd(host string) *SSHCmd {
	ssh := &SSHCmd{
		Host: host,
	}

	ssh.Cmd = NewCmd()
	ssh.Cmd.Prefix.stdout = color(host) + " "
	ssh.Cmd.Prefix.stderr = color(host) + colorErr("@err ")
	ssh.Cmd.Prefix.cmd = "$ "
	return ssh
}

func (s *SSHCmd) warpInSsh(command string) (sshArgs []string) {
	sshArgs = append(sshArgs, s.SshBinPath)

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

func (s *SSHCmd) Wait() error {
	return s.Cmd.Wait()
}

func (s *SSHCmd) Run(command string) (res CmdRes, err error) {
	if res, err = s.Start(command); err != nil {
		return
	}

	err = s.Wait()

	return
}

func (s *SSHCmd) Start(command string) (res CmdRes, err error) {
	if s.SshBinPath == "" {
		if s.SshBinPath, err = FindPath(sshBinList); err != nil {
			err = errors.Wrapf(err, "can't find ssh binary: %v", shellPathList)
			return
		}
	}

	sshArgs := s.warpInSsh(command)

	res, err = s.Cmd.Start(strings.Join(sshArgs, " "))

	return
}
