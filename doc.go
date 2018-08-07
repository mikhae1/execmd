// Copyright(c) 2018 by Mink0. All rights reserved.

/*
Package execmd is a Golang library providing a simple interface to shell commands execution

Features
		- execute commands in system shell
		- you could use shell variables, pipes, redirections
		- execute remote shell commands
		- interface is based on os/exec
		- realtime `stdout` and `stderr` output with fancy colors and prefixes
		- remote commands execution is implemented by wrapping standart OpenSSH client
		- all your ssh configuration (including ssh agent forwarding) works
		- parallel and serial remote command execution supported

Documentation
	See https://github.com/mink0/exec-cmd/blob/master/README.md for full documentation
*/
package execmd
