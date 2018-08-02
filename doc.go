// Copyright 2018 by Mink0. All rights reserved.

/*
Exec-cmd is a Golang library providing a simple interface to shell commands execution

Features
* execute command in system shell, so you could use variables, pipes, redirections
* based on [exec](https://golang.org/pkg/os/exec/)
* execute local and remote shell commands
* remote commands execution is implemented by wrapping standart [OpenSSH](https://www.openssh.com/) SSH client, so all your ssh configuration (including ssh agent forwarding) works as expected
* realtime `stdout` and `stderr` output with fancy colors and prefixes
* run remote commands on several hosts (parallel and serial execution supported)
*/
package execmd
