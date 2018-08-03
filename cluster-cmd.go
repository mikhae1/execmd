package execmd

// ClusterSSHCmd is a wrapper on SSHCmd
type ClusterSSHCmd struct {
	StopOnError bool
	Hosts       []string
	SSHCmds     []*SSHCmd

	StartedCmds []ClusterRes
}

// ClusterRes contains resultss of command execution:
// res - stdout, stderr
// err - error
type ClusterRes struct {
	Host string
	Res  CmdRes
	Err  error
}

// NewClusterSSHCmd init ClusterSSHCmd with defaults
func NewClusterSSHCmd(hosts []string) (c ClusterSSHCmd) {
	c = ClusterSSHCmd{}
	c.StopOnError = false
	c.Hosts = append([]string(nil), hosts...)
	for _, h := range hosts {
		c.SSHCmds = append(c.SSHCmds, NewSSHCmd(h))
	}

	return
}

// Wait wraps SSHCmd.Wait for array of hosts into c.StartedCmds struct
func (c *ClusterSSHCmd) Wait() error {
	// TODO: list errors should be returned, or maybe hostname appended
	// now you should access `.StartedCmds` to see exact where error occurs
	var lastError error
	for i := range c.StartedCmds {
		// skip errors on Start()
		if c.StartedCmds[i].Err != nil {
			continue
		}

		c.StartedCmds[i].Err = c.SSHCmds[i].Wait()
		if c.StartedCmds[i].Err != nil {
			if c.StopOnError {
				return c.StartedCmds[i].Err
			}

			lastError = c.StartedCmds[i].Err
		}
	}

	return lastError
}

// Run executes command in parallel: all commands starts running simultaniosly at the hosts
func (c *ClusterSSHCmd) Run(command string) (results []ClusterRes, err error) {
	if results, err = c.Start(command); err != nil {
		return
	}

	err = c.Wait()

	// FIXME
	// c.StartedCmds results with host filled
	// protect the results from changes from ClusterSSHCmd
	results = append([]ClusterRes(nil), c.StartedCmds...)
	return
}

// RunOneByOne runs command in series: run at first host, then run at second, then...
func (c *ClusterSSHCmd) RunOneByOne(command string) (results []ClusterRes, err error) {
	return c.startAndRun(command, false)
}

// Start runs command in parallel
func (c *ClusterSSHCmd) Start(command string) (results []ClusterRes, err error) {
	return c.startAndRun(command, true)
}

// Loop through hosts and
// .Start() or .Run() ssh command depending on `start` flag
func (c *ClusterSSHCmd) startAndRun(command string, start bool) (results []ClusterRes, err error) {
	// reset started on each new start
	c.StartedCmds = []ClusterRes{}
	for i, host := range c.Hosts {
		cres := ClusterRes{}
		cres.Host = host
		if start {
			cres.Res, cres.Err = c.SSHCmds[i].Start(command)
		} else {
			cres.Res, cres.Err = c.SSHCmds[i].Run(command)
		}
		results = append(results, cres)

		// save results
		c.StartedCmds = append(c.StartedCmds, cres)
		if c.StopOnError && cres.Err != nil {
			err = cres.Err
			return
		}
	}

	return
}
