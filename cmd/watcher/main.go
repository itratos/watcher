package main

import (
	"flag"
	"fmt"
	"github.com/itratos/watcher/watcher"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
	"unicode"
)

func main() {
	flags := flag.NewFlagSet("watcher", flag.ExitOnError)

	verbose := flags.Bool("verbose", false, "verbose event notifications")
	interval := flags.String("interval", "100ms", "watcher poll interval")
	recursive := flags.Bool("recursive", false, "watch directories recursively")
	allfiles := flags.Bool("all", false, "watch all files, including dotfiles")
	cmd := flags.String("cmd", "", "command to run when an event occurs")
	startcmd := flags.Bool("startcmd", false, "run the command when watcher starts")
	listFiles := flags.Bool("list", false, "list watched files on start")
	stdinPipe := flags.Bool("pipe", false, "pipe event's info to command's stdin")
	keepalive := flags.Bool("keepalive", false, "keep alive when a cmd returns code != 0")
	ignore := flags.String("ignore", "", "comma separated list of paths to ignore")
	version := flags.Bool("version", false, "prints current version")
	help := flags.Bool("help", false, "prints this help")

	flags.BoolVar(recursive, "r", false, "watch directories recursively")
	flags.BoolVar(allfiles, "a", false, "watch all files, including dotfiles")
	flags.BoolVar(version, "v", false, "prints current version")
	flags.BoolVar(help, "h", false, "prints this help")

	flags.Usage = func() {
		fmt.Fprintf(flags.Output(), "USAGE:\n  %s [OPTIONS] [PATHS]\n\nOPTIONS:\n", flags.Name())
		flags.PrintDefaults()
	}

	flags.Parse(os.Args[1:])

	const CmdVersion = "v2.0.4"

	if *version {
		fmt.Println(CmdVersion)
		os.Exit(0)
	}

	if *help {
		flags.SetOutput(os.Stdout)
		flags.Usage()
		os.Exit(0)
	}

	// Retrieve the list of files and folders.
	files := flags.Args()

	// If no files/folders were specified, watch the current directory.
	if len(files) == 0 {
		curDir, err := os.Getwd()
		if err != nil {
			log.Fatalln(err)
		}
		files = append(files, curDir)
	}

	var cmdName string
	var cmdArgs []string
	if *cmd != "" {
		split := strings.FieldsFunc(*cmd, unicode.IsSpace)
		cmdName = split[0]
		if len(split) > 1 {
			cmdArgs = split[1:]
		}
	}

	// Create a new Watcher with the specified options.
	w := watcher.New()
	w.IgnoreHiddenFiles(!*allfiles)

	// Get any of the paths to ignore.
	ignoredPaths := strings.Split(*ignore, ",")

	for _, path := range ignoredPaths {
		trimmed := strings.TrimSpace(path)
		if trimmed == "" {
			continue
		}

		err := w.Ignore(trimmed)
		if err != nil {
			log.Fatalln(err)
		}
	}

	done := make(chan struct{})
	go func() {
		defer close(done)

		for {
			select {
			case event := <-w.Event:
				var eventStr string
				if *verbose {
					eventStr = event.VerboseString()
				} else {
					eventStr = event.String()
				}

				// Print the event's info.
				fmt.Println(eventStr)

				// Run the command if one was specified.
				if *cmd != "" {
					c := exec.Command(cmdName, cmdArgs...)
					if *stdinPipe {
						c.Stdin = strings.NewReader(eventStr)
					} else {
						c.Stdin = os.Stdin
					}
					c.Stdout = os.Stdout
					c.Stderr = os.Stderr
					if err := c.Run(); err != nil {
						if (c.ProcessState == nil || !c.ProcessState.Success()) && *keepalive {
							log.Println(err)
							continue
						}
						log.Fatalln(err)
					}
				}
			case err := <-w.Error:
				if err == watcher.ErrWatchedFileDeleted {
					fmt.Println(err)
					continue
				}
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Add the files and folders specified.
	for _, file := range files {
		if *recursive {
			if err := w.AddRecursive(file); err != nil {
				log.Fatalln(err)
			}
		} else {
			if err := w.Add(file); err != nil {
				log.Fatalln(err)
			}
		}
	}

	// Print a list of all of the files and folders being watched.
	if *listFiles {
		for path, f := range w.WatchedFiles() {
			fmt.Printf("%s: %s\n", path, f.Name())
		}
		fmt.Println()
	}

	fmt.Printf("Watching %d files\n", len(w.WatchedFiles()))

	// Parse the interval string into a time.Duration.
	parsedInterval, err := time.ParseDuration(*interval)
	if err != nil {
		log.Fatalln(err)
	}

	closed := make(chan struct{})

	c := make(chan os.Signal)
	signal.Notify(c, os.Kill, os.Interrupt)
	go func() {
		<-c
		w.Close()
		<-done
		fmt.Println("watcher closed")
		close(closed)
	}()

	// Run the command before watcher starts if one was specified.
	go func() {
		if *cmd != "" && *startcmd {
			c := exec.Command(cmdName, cmdArgs...)
			c.Stdin = os.Stdin
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			if err := c.Run(); err != nil {
				if (c.ProcessState == nil || !c.ProcessState.Success()) && *keepalive {
					log.Println(err)
					return
				}
				log.Fatalln(err)
			}
		}
	}()

	// Start the watching process.
	if err := w.Start(parsedInterval); err != nil {
		log.Fatalln(err)
	}

	<-closed
}
