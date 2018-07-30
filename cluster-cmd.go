package execmd

type ClusterSSHCmd struct {
	StopOnError bool
	Hosts       []string
	SSHCmds     []*SSHCmd
	Started     []ClusterRes
}

type ClusterRes struct {
	res CmdRes
	err error
}

func NewClusterSSHCmd(hosts []string) (c ClusterSSHCmd) {
	c = ClusterSSHCmd{}
	c.StopOnError = false
	c.Hosts = append([]string(nil), hosts...)
	for _, h := range hosts {
		c.SSHCmds = append(c.SSHCmds, NewSSHCmd(h))
	}

	return
}

func (c *ClusterSSHCmd) Wait() error {
	// TODO: list errors should be returned, or maybe hostname appended
	// now you should access `.Started` to see where exact error occur
	var lastError error
	for i := range c.Started {
		// skip errors on Start()
		if c.Started[i].err != nil {
			continue
		}

		c.Started[i].err = c.SSHCmds[i].Wait()
		if c.Started[i].err != nil {
			if c.StopOnError {
				return c.Started[i].err
			}

			lastError = c.Started[i].err
		}
	}

	return lastError
}

// Run command in parallel: all commands starts running simultaniosly at the hosts
func (c *ClusterSSHCmd) Run(command string) (cresult []ClusterRes, err error) {
	if cresult, err = c.Start(command); err != nil {
		return
	}

	err = c.Wait()

	cresult = append([]ClusterRes(nil), c.Started...)
	return
}

// RunSeq Run command in series: run at first host, then run at second, then...
func (c *ClusterSSHCmd) RunSeq(command string) (cresult []ClusterRes, err error) {
	for i := range c.Hosts {
		cres := ClusterRes{}
		cres.res, cres.err = c.SSHCmds[i].Run(command)
		cresult = append(cresult, cres)

		if c.StopOnError && cres.err != nil {
			err = cres.err
			return
		}
	}

	return
}

// Start starts command in parallel
func (c *ClusterSSHCmd) Start(command string) (cresult []ClusterRes, err error) {
	// reset started on each start
	c.Started = []ClusterRes{}

	for i := range c.Hosts {
		cres := ClusterRes{}
		cres.res, cres.err = c.SSHCmds[i].Start(command)
		cresult = append(cresult, cres)

		// save result
		c.Started = append(c.Started, cres)
		if c.StopOnError && cres.err != nil {
			err = cres.err
			return
		}
	}

	return
}
