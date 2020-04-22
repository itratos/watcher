# watcher

Forked from https://github.com/radovskyb/watcher to pull in updates and modify for personal usage.

## Installation

```shell
go get -u github.com/dideler/watcher/...
```

## Executable

The `watcher` CLI command is installed when using the `go get` command from above.

```
USAGE:
  watcher [OPTIONS] [PATHS]

OPTIONS:
  -a
  -all
    	watch all files, including dotfiles
  -cmd string
    	command to run when an event occurs
  -h
  -help
    	prints this help
  -ignore string
    	comma separated list of paths to ignore
  -interval string
    	watcher poll interval (default "100ms")
  -keepalive
    	keep alive when a cmd returns code != 0
  -list
    	list watched files on start
  -pipe
    	pipe event's info to command's stdin
  -r
  -recursive
    	watch directories recursively
  -startcmd
    	run the command when watcher starts
  -v
  -version
    	prints current version
```

Without any arguments, it watches the current directory for changes and notifies on any events that occur.
