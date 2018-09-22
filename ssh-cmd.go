/*
	Wrapper for Cmd to invoke ssh commands via OpenSSH binary

	Copyright(c) 2018 mink0
*/

package execmd

import (
	"os"
	"strings"

	"github.com/pkg/errors"
)

// SSHCmd is a wrapper on Cmd
type SSHCmd struct {
	Cmd *Cmd

	Interactive   bool
	SSHExecutable string
	Host          string
	User          string
	Port          string
	KeyPath       string
	Cwd           string
}

// NewSSHCmd initializes SSHCmd with defaults
func NewSSHCmd(host string) *SSHCmd {
	ssh := SSHCmd{
		Host:          host,
		SSHExecutable: "ssh",
	}

	ssh.Cmd = NewCmd()
	ssh.Cmd.PrefixStdout = color(host) + " "
	ssh.Cmd.PrefixStderr = color(host) + colorErr("@err ")

	// Path to ssh binary could be overridden by setting `SSH_EXECUTABLE` env variable
	if sshEnvExec, ok := os.LookupEnv("SSH_EXECUTABLE"); ok {
		ssh.SSHExecutable = sshEnvExec
	}

	// User detect from user@host
	if arr := strings.Split(host, "@"); len(arr) == 2 {
		ssh.User = arr[0]
		ssh.Host = arr[1]
	}

	return &ssh
}

// Wait wraps Cmd.Wait()
func (s *SSHCmd) Wait() error {
	return s.Cmd.Wait()
}

// Run wraps Cmd.Run()
func (s *SSHCmd) Run(command string) (res CmdRes, err error) {
	if res, err = s.Start(command); err != nil {
		return
	}

	err = s.Wait()
	return
}

// Start wraps Cmd.Start() with ssh invocation
func (s *SSHCmd) Start(command string) (res CmdRes, err error) {
	if s.Host == "" {
		err = errors.New("no host to run ssh command")
		return
	}

	sshArgs := s.warpInSSH(command)

	res, err = s.Cmd.Start(strings.Join(sshArgs, " "))

	return
}

// transform `command` into ssh-compatible argument string
func (s *SSHCmd) warpInSSH(command string) (sshArgs []string) {
	sshArgs = append(sshArgs, s.SSHExecutable)

	hostWithUser := s.Host
	if s.User != "" {
		hostWithUser = s.User + "@" + s.Host
	}

	sshArgs = append(sshArgs, hostWithUser)

	if s.Interactive || strings.Contains(command, "sudo") {
		sshArgs = append(sshArgs, "-tt")
	}
	if s.Port != "" {
		sshArgs = append(sshArgs, "-p", s.Port)
	}
	if s.KeyPath != "" {
		sshArgs = append(sshArgs, "-i", s.KeyPath)
	}
	if s.Cwd != "" {
		command = "cd " + s.Cwd + " && " + command
	}

	// escape single quotes for shell encapsulation
	sshArgs = append(sshArgs, "'"+strings.Replace(command, "'", "'\\''", -1)+"'")

	return
}
