package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func run(e env, c config) {
	if len(c.cmd) == 0 {
		os.Exit(0)
	}

	exe, err := exec.LookPath(c.cmd[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %s: %s\n", c.prog, err)
		os.Exit(1)
	}

	err = syscall.Exec(exe, c.cmd, e.Environ())
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %s: %s\n", c.prog, err)
		os.Exit(1)
	}
}
