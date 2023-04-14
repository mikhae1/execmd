# execmd

`ExecCmd` is a user-friendly Go package that offers a simplified interface for shell command execution. Built on top of the [exec](https://golang.org/pkg/os/exec/) package, `ExecCmd` enables command invocation in a system shell and combines multiple stdout and stderr into a single stdout with prefixes. It supports both local and remote command execution, with remote commands implemented through the [OpenSSH](https://www.openssh.com/) binary.

## Key features

- Easy way to execute local and remote shell commands
- Support for shell environment variables, pipes, and redirection
- Compatibility with system SSH configuration (including ssh-agent forwarding)
- Run commands on multiple remote hosts (ideal for cluster operations) with parallel or serial execution options
- Real-time `stdout` and `stderr` output featuring auto coloring and prefixing
- Output buffers can be captured for programmatic access
- Minimum number of third party dependencies

## Installation

    go get "github.com/mink0/exec-cmd"

Then import `exec-cmd` in your application:

```go
import "github.com/mink0/exec-cmd"
```

## Examples

### Local command execution

```go
package main

import "github.com/mink0/exec-cmd"

func main() {
  // run local command in a shell
  execmd.NewCmd().Run("ps aux | grep go")
}
```

### Remote command execution

```go
package main

import "github.com/mink0/exec-cmd"

func main() {
  // run command on a remote host using ssh
  remote := execmd.NewSSHCmd("192.168.1.194")
  res, err := remote.Run(`VAR="$(hostname)"; echo "hello $VAR"`)

  if err == nil {
    fmt.Printf("captured output: %s", res.Stdout)
  }
}
```

Results:
```sh
$ /usr/bin/ssh 192.168.1.194 'VAR="$(hostname)"; echo "hello $VAR"'
192.168.1.194 hello host-01.local
captured output: hello host-01.local
```

### Remote cluster command execution

```go
cluster := execmd.NewClusterSSHCmd([]string{"host-01", "host-02", "host-03"})

// execute in parallel order
res, err := cluster.Run(`VAR=std; echo "Hello $VAR out"; echo Hello $VAR err >&2`)

// execute in serial order
res, err = cluster.RunOneByOne(`VAR=std; echo "Hello $VAR out"`)
}
```

Parallel execution results:
```sh
$ /usr/bin/ssh host-01 'VAR=std; echo "Hello $VAR out"; echo "Hello $VAR err" >&2'
$ /usr/bin/ssh host-02 'VAR=std; echo "Hello $VAR out"; echo "Hello $VAR err" >&2'
$ /usr/bin/ssh host-03 'VAR=std; echo "Hello $VAR out"; echo "Hello $VAR err" >&2'
host-01 Hello std out
host-01@err Hello std err
host-03@err Hello std err
host-03 Hello std out
host-02 Hello std out
host-02@err Hello std err
```

Serial execution results:
```sh
$ /usr/bin/ssh host-01 'VAR=std; echo "Hello $VAR out"; echo Hello $VAR err >&2'
host-01 Hello std out
host-01@err Hello std err
$ /usr/bin/ssh host-02 'VAR=std; echo "Hello $VAR out"; echo Hello $VAR err >&2'
host-02@err Hello std err
host-02 Hello std out
$ /usr/bin/ssh host-03 'VAR=std; echo "Hello $VAR out"; echo Hello $VAR err >&2'
host-03@err Hello std err
host-03 Hello std out
```

## Testing

You should enable `SSH` server locally and add your personal ssh key to `known_hosts` to avoid password prompting:

```sh
ssh-copy-id 127.0.0.1
ssh-copy-id localhost
```

Run tests:

    go test

## License

The MIT License (MIT) - see [LICENSE](./LICENSE) for details.