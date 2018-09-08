package execmd

// ClusterSSHCmd is a wrapper on SSHCmd
type ClusterSSHCmd struct {
	Cmds []ClusterCmd

	StopOnError bool
}

// ClusterCmd wraps SSHCmd and preserves host name and saves error from .start() for .Wait() method
type ClusterCmd struct {
	SSHCmd SSHCmd

	Host string
}

// ClusterRes contains results of command execution
type ClusterRes struct {
	Host string
	Err  error
	Res  CmdRes
}

// NewClusterSSHCmd inits ClusterSSHCmd with defaults
func NewClusterSSHCmd(hosts []string) *ClusterSSHCmd {
	c := ClusterSSHCmd{}
	c.StopOnError = false
	c.Cmds = make([]ClusterCmd, len(hosts))
	for i, host := range hosts {
		c.Cmds[i].Host = host
		c.Cmds[i].SSHCmd = *NewSSHCmd(host)
	}
	return &c
}

// Loop through the hosts and run .Start() or .Run() method (depend on `parallel` flag)
func (c *ClusterSSHCmd) start(command string, parallel bool) ([]ClusterRes, error) {

	results := make([]ClusterRes, len(c.Cmds))
	for i, cmd := range c.Cmds {
		results[i].Host = cmd.Host

		// no need to implement interfaces here, we always have only: .Start() and .Run() methods
		exec := cmd.SSHCmd.Start
		if !parallel {
			exec = cmd.SSHCmd.Run
		}

		results[i].Res, results[i].Err = exec(command)

		if c.StopOnError && results[i].Err != nil {
			return results[:i+1], results[i].Err
		}
	}

	return results, nil
}

// Wait calls SSHCmd.Wait for the list of Cmds []ClusterCmd
// you should access `.Cmds` to see exact where error occurs
func (c *ClusterSSHCmd) Wait() (err error) {
	for _, cmd := range c.Cmds {

		err = cmd.SSHCmd.Wait()
		if c.StopOnError && err != nil {
			return
		}
	}

	return
}

// Run executes command in parallel and waits for the results. Command starts simultaneously at each of the hosts.
func (c *ClusterSSHCmd) Run(command string) (results []ClusterRes, err error) {
	if results, err = c.Start(command); err != nil {
		return
	}

	err = c.Wait()

	return
}

// RunOneByOne executes command in series: run at first host, then run at second host, then...
func (c *ClusterSSHCmd) RunOneByOne(command string) (results []ClusterRes, err error) {
	return c.start(command, false)
}

// Start executes command in parallel, no wait for the results
func (c *ClusterSSHCmd) Start(command string) (results []ClusterRes, err error) {
	return c.start(command, true)
}
