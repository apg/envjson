package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
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
