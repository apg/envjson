// Program envjson provides a mechanism for managing an environment before
// launching a program. It uses a JSON spec to ensure that all of the
// expected variables are in the environment (reading from either STDIN, or
// the actual environment itself).
//
// See the README for more information.
//
// Released under the Modified BSD license. See LICENSE for more information.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

func forwardSignals(cmd *exec.Cmd) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc)
	go func() {
		for {
			s := <-sigc
			cmd.Process.Signal(s)
		}
	}()
}

func run(e env, c config) {
	if len(c.cmd) == 0 {
		os.Exit(0)
	}

	cmd := exec.Command(c.cmd[0], c.cmd[1:]...)
	cmd.Env = e.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: fatal: failed to run command: %s\n", c.prog, err)
		os.Exit(1)
	}

	forwardSignals(cmd)

	err = cmd.Wait()
	if err == nil {
		os.Exit(0) // Successful.
	} else if exit, ok := err.(*exec.ExitError); ok {
		if s, ok := exit.Sys().(syscall.WaitStatus); ok {
			os.Exit(s.ExitStatus())
		}
	}

	fmt.Fprintf(os.Stderr, "%s: fatal: failed to run command: %s\n", os.Args[0], err)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [--stdin | --help] JSONfile <cmd> <cmd args>\n", os.Args[0])
}

type config struct {
	displayHelp bool
	displayEnv  bool
	fromStdin   bool
	cleanEnv    bool
	prog        string
	specFile    string
	cmd         []string
}

func configure(args []string) (c config, err error) {
	c.prog = args[0]

	args = args[1:]

	if len(args) == 0 { // by default just display the environment
		c.displayEnv = true
		return
	}

	i := 0

loop:
	for ; i < len(args); i++ {
		switch {
		case os.Args[i] == "--help":
			c.displayHelp = true
			return
		case args[i] == "--stdin" || args[i] == "-":
			c.fromStdin = true
		case args[i] == "-i" || args[i] == "--ignore-environment":
			c.cleanEnv = true
		case strings.HasPrefix(args[i], "--") || strings.HasPrefix(args[i], "-"):
			err = fmt.Errorf("invalid argument %q", args[i])
			return
		default:
			break loop
		}
	}

	if i < len(args) {
		c.specFile = args[i] // save off the spec file for loading later.

		i++
		if i < len(args) {
			c.cmd = args[i:]
		} else {
			c.displayEnv = true
		}
	} else {
		c.displayEnv = true
	}

	return
}

func main() {
	config, err := configure(os.Args)
	if config.displayHelp {
		usage()
		os.Exit(1)
	} else if err != nil {
		fmt.Fprintf(os.Stderr,
			"%s: error: %s\n",
			config.prog, err)
		usage()
		os.Exit(1)
	}

	// The config we *want* to launch command with
	local := make(env)

	// The config we originally *have*
	parent := make(env)

	if !config.cleanEnv {
		parent.FromEnviron(os.Environ())
	}

	if config.fromStdin {
		err := parent.Read(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr,
				"%s: error: Unable to read parent env from STDIN: %s\n",
				config.prog, err)
			os.Exit(1)
		}
	}

	// Is there a spec file? If so, read it.
	if config.specFile != "" {
		err := local.FromFile(config.specFile)
		if err != nil {
			fmt.Fprintf(os.Stderr,
				"%s: error: Unable to read env from %q: %s\n",
				config.prog, config.specFile, err)
			os.Exit(1)
		}
	}

	// Validate that we have what we need.
	err = local.Merge(parent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: error: Environment check failed: %s\n",
			config.prog, err)
		os.Exit(1)
	}

	// We're just dumping to the screen.
	if config.displayEnv {
		err = local.Dump(os.Stdout)
		if err != nil {
			fmt.Fprintf(os.Stderr,
				"%s: error: Unable to dump env to JSON: %s\n",
				config.prog, err)
			os.Exit(1)
		}
		return
	}

	run(local, config)
}
