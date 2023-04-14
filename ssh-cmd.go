package execmd

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// SSHCmd is a wrapper on Cmd to invoke ssh commands via OpenSSH binary
type SSHCmd struct {
	Cmd           *Cmd
	Interactive   bool
	SSHExecutable string
	Host          string
	User          string
	Port          string
	KeyPath       string
	Cwd           string
}

// NewSSHCmd initializes SSHCmd with defaults and sets the target host
func NewSSHCmd(host string) *SSHCmd {
	ssh := &SSHCmd{
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

	return ssh
}

// Wait wraps Cmd.Wait(), waiting for the remote command to complete
func (s *SSHCmd) Wait() error {
	return s.Cmd.Wait()
}

// Run wraps Cmd.Run(), executing the remote command and waiting for it to complete
func (s *SSHCmd) Run(command string, timeout ...time.Duration) (res CmdRes, err error) {
	if res, err = s.Start(command, timeout...); err != nil {
		return
	}

	err = s.Wait()
	return
}

// Start wraps Cmd.Start() with ssh invocation, starting the remote command
func (s *SSHCmd) Start(command string, timeout ...time.Duration) (res CmdRes, err error) {
	if s.Host == "" {
		err = fmt.Errorf("no host to run ssh command")
		return
	}

	sshArgs, err := s.warpInSSH(command)
	if err != nil {
		return res, fmt.Errorf("failed to prepare ssh command: %w", err)
	}

	res, err = s.Cmd.Start(strings.Join(sshArgs, " "), timeout...)
	return
}

// warpInSSH takes a command string and returns an ssh-compatible argument slice
func (s *SSHCmd) warpInSSH(command string) ([]string, error) {
	sshArgs := []string{s.SSHExecutable}
	hostWithUser := s.Host
	if s.User != "" {
		hostWithUser = s.User + "@" + s.Host
	}

	sshArgs = append(sshArgs, hostWithUser)

	if s.Interactive || strings.Contains(command, "sudo") {
		sshArgs = append(sshArgs, "-tt")
		s.Cmd.Interactive = true
	}
	if s.Port != "" {
		sshArgs = append(sshArgs, "-p", s.Port)
	}
	if s.KeyPath != "" {
		if _, err := os.Stat(s.KeyPath); err != nil {
			return nil, fmt.Errorf("ssh key not found at path %s: %w", s.KeyPath, err)
		}
		sshArgs = append(sshArgs, "-i", s.KeyPath)
	}
	if s.Cwd != "" {
		command = "cd " + s.Cwd + " && " + command
	}

	// escape single quotes for shell encapsulation
	sshArgs = append(sshArgs, "'"+strings.Replace(command, "'", "'\\''", -1)+"'")
	return sshArgs, nil
}
