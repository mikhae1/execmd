package execmd

import (
	"fmt"
	"time"
)

// ClusterSSHCmd is a wrapper on SSHCmd that allows executing commands on multiple hosts in parallel or sequentially.
type ClusterSSHCmd struct {
	Cmds   []ClusterCmd
	Errors []error

	Cwd         string
	StopOnError bool
}

// ClusterCmd wraps SSHCmd and preserves the host name, and saves errors from .Start() for the .Wait() method.
type ClusterCmd struct {
	SSHCmd SSHCmd

	Host string
}

// ClusterRes contains the results of the command execution.
type ClusterRes struct {
	Host string
	Err  error
	Res  CmdRes
}

// NewClusterSSHCmd initializes ClusterSSHCmd with defaults.
func NewClusterSSHCmd(hosts []string) *ClusterSSHCmd {
	c := ClusterSSHCmd{}
	c.StopOnError = false
	c.Cmds = make([]ClusterCmd, len(hosts))
	c.Errors = make([]error, len(hosts))
	for i, host := range hosts {
		c.Cmds[i].Host = host
		c.Cmds[i].SSHCmd = *NewSSHCmd(host)
	}
	return &c
}

// start iterates through the hosts and runs .Start() or .Run() method (depends on `parallel` flag).
func (c *ClusterSSHCmd) start(command string, parallel bool, timeout ...time.Duration) ([]ClusterRes, error) {
	results := make([]ClusterRes, len(c.Cmds))
	for i, cmd := range c.Cmds {
		// Set cluster common variables
		if c.Cwd != "" {
			cmd.SSHCmd.Cwd = c.Cwd
		}

		results[i].Host = cmd.Host

		// No need to implement full interfaces here, we use only: .Start() and .Run() methods
		exec := cmd.SSHCmd.Start
		if !parallel {
			exec = cmd.SSHCmd.Run
		}

		results[i].Res, results[i].Err = exec(command, timeout...)

		if c.StopOnError && results[i].Err != nil {
			return results[:i+1], fmt.Errorf("error on host %s: %w", cmd.Host, results[i].Err)
		}
	}

	return results, nil
}

// Wait calls SSHCmd.Wait for each Cmd in the list of ClusterCmds.
// It returns the first caught .Wait() error ans stops if .StopOnError is true.
// To see underlying SSHCmd command errors, access the .Cmds attribute.
func (c *ClusterSSHCmd) Wait() error {
	var firstErr error
	for i, cmd := range c.Cmds {
		if err := cmd.SSHCmd.Wait(); err != nil {
			c.Errors[i] = err

			if firstErr == nil {
				firstErr = fmt.Errorf("error on host %s: %w", cmd.Host, err)
				if c.StopOnError {
					return firstErr
				}
			}
		}
	}

	return firstErr
}

// Run executes a command in parallel on all hosts and waits for the results.
// The command starts simultaneously on each host.
// It returns results and the first caught error.
// To see underlying SSHCmd command errors, access the .Cmds attribute.
func (c *ClusterSSHCmd) Run(command string, timeout ...time.Duration) (results []ClusterRes, err error) {
	if results, err = c.Start(command, timeout...); err != nil {
		return
	}

	err = c.Wait()

	// populate .Err
	for i, err := range c.Errors {
		results[i].Err = err
	}

	return
}

// RunOneByOne executes a command in series: run at the first host, then run at the second host, and so on.
// It returns results and the first caught error.
// To see underlying SSHCmd command errors, access the .Cmds attribute.
func (c *ClusterSSHCmd) RunOneByOne(command string, timeout ...time.Duration) (results []ClusterRes, err error) {
	return c.start(command, false, timeout...)
}

// Start executes a command in parallel on all hosts without waiting for the results.
// The command starts simultaneously on each host.
// It returns results and the first caught error.
func (c *ClusterSSHCmd) Start(command string, timeout ...time.Duration) (results []ClusterRes, err error) {
	return c.start(command, true, timeout...)
}
