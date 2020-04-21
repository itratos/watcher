# watcher

Forked from https://github.com/radovskyb/watcher to pull in updates and modify for personal usage.

## Installation

```shell
go get -u github.com/dideler/watcher
```

## Executable

`watcher` comes with a CLI command which is installed when using the `go get` command from above.

```
Usage: watcher [OPTS] [FILES]

Options:

  -cmd string
    	command to run when an event occurs
  -dotfiles
    	watch dot files (default true)
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
  -recursive
    	watch folders recursively (default true)
  -startcmd
    	run the command when watcher starts
```

Without any arguments, it watches the current directory recursively for changes and notifies any events that occur.