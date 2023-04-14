/*
Package execmd is a Golang library providing a simple interface to shell commands execution

Features
  - Execute commands in the system shell.
  - Utilize shell variables, pipes, and redirections.
  - Execute remote shell commands.
  - Interface is based on the os/exec package.
  - Real-time stdout and stderr output with customizable colors and prefixes.
  - Remote commands execution is implemented using the standard OpenSSH client.
  - Supports all SSH configurations, including SSH agent forwarding.
  - Parallel and serial remote command execution supported.

Documentation

	See https://github.com/mikhae1/exec-cmd/blob/master/README.md for full documentation
*/
package execmd
