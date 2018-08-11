# Exec-cmd

Exec-cmd is a Golang library providing a simple interface to shell commands execution

## Key features

* execute command in system shell, so you could use variables, pipes, redirections
* execute local and remote shell commands
* interface is similar to [exec](https://golang.org/pkg/os/exec/)
* real time `stdout` and `stderr` output with fancy colors and prefixes
* remote commands execution is implemented by wrapping [OpenSSH](https://www.openssh.com/) SSH client,
  so all your ssh configuration (including ssh agent forwarding) works as expected
* run remote commands on several hosts (parallel and serial execution supported)

## Installation

Install the latest version of the library:

* using `go get`:

      go get -u "github.com/mink0/exec-cmd"

* using [golang/dep](https://github.com/golang/dep) tool:

      dep ensure -add "github.com/mink0/exec-cmd"

Include `exec-cmd` in your application:

```go
import "github.com/mink0/exec-cmd"
```

## Local command example

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

## Remote command example

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

## Remote cluster command example

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
