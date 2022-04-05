# execmd

`execmd` is a simple Go package providing an interface for shell command execution.
It wraps [exec](https://golang.org/pkg/os/exec/) to invoke command in a system shell and redirects multiple `stdout`, `stderr` into a single `stdout` using prefixes. It supports both local and remote commands execution, remote commands are implemented by invoking [OpenSSH](https://www.openssh.com/) binary.

## Key features

* local and remote commands execution
* simple interface similar to [exec](https://golang.org/pkg/os/exec/)
* shell environment variables, pipes and redirection are supported
* system ssh configuration is supported (including `ssh-agent` forwarding)
* you can run single command on multiple hosts (for cluster operations): parallel and serial execution is supported
* real time `stdout` and `stderr` output with fancy coloring and prefixing
* output buffers could be captured to access them programmatically

## Installation

* using go modules:

      go mod vendor

* using `go get`:

      go get -u "github.com/mink0/exec-cmd"

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
  execmd.NewCmd().Run("ps aux | grep go")
}
```

```sh
$ ps aux | grep go
> mink0         9203   0.6  0.0 558434044   2220 s001  S+    3:11PM   0:00.01 ./go-dp
> mink0         1934   0.0  0.2 558459904  15744   ??  S    11:20AM   0:01.71 /Users/mink0/go/bin/gocode -s -sock unix -addr 127.0.0.1:37373
```

### Remote command execution

```go
package main

import "github.com/mink0/exec-cmd"

func main() {
  remote := execmd.NewSSHCmd("192.168.1.194")
  res, err := remote.Run(`VAR="$(hostname)"; echo "hello $VAR"`)
  if err == nil {
    fmt.Printf("saved output: %s", res.Stdout)
  }
}
```

```sh
$ /usr/bin/ssh 192.168.1.194 'VAR="$(hostname)"; echo "hello $VAR"'
192.168.1.194 hello host-01.local
saved output: hello host-01.local
```

### Remote cluster command execution

```go
package main

import "github.com/mink0/exec-cmd"

cluster := execmd.NewClusterSSHCmd([]string{"host-01", "host-02", "host-03"})

func main() {
  // execute in parallel order
  res, err := cluster.Run(`VAR=std; echo "Hello $VAR out"; echo Hello $VAR err >&2`)
  if err == nil {
    fmt.Printf("saved output: %v", res)
  }

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
