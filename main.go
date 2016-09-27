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
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

var verbose = true

type config map[string]configValue

type configValue struct {
	Value    string
	Required bool
	Inherit  bool

	isSet bool
}

func (c *configValue) UnmarshalJSON(b []byte) error {
	if b[0] == '"' && len(b) >= 2 {
		c.Value = string(b[1 : len(b)-1])
		c.isSet = true
		return nil
	}

	var target struct {
		Value    string `json:"value"`
		Required bool   `json:"required"`
		Inherit  bool   `json:"inherit"`
		Doc      string `json:"doc"`
	}

	err := json.Unmarshal(b, &target)
	if err != nil {
		return err
	}

	if target.Value != "" {
		c.isSet = true
	}

	c.Value = target.Value
	c.Required = target.Required
	c.Inherit = target.Inherit

	return nil
}

// Read populates a config from a reader.
func (c config) Read(in io.Reader) error {
	dec := json.NewDecoder(in)
	err := dec.Decode(&c)

	return err
}

// FromFile reads in configValues from a file
func (c config) FromFile(f string) error {
	in, err := os.Open(f)
	if err != nil {
		return err
	}

	return c.Read(in)

}

// FromEnviron reads in configValues from the environment
func (c config) FromEnviron() error {
	for _, kv := range os.Environ() {
		bits := strings.SplitN(kv, "=", 2)
		if len(bits) == 2 {
			c[bits[0]] = configValue{Value: bits[1], isSet: true}
		} else {
			c[bits[0]] = configValue{Value: ""}
		}
	}
	return nil
}

// Inherit attempts to find name in `p`
func (c config) Inherit(name string, p config) error {
	pv, ok := p[name]
	if !ok {
		return fmt.Errorf("Unable to inherit %q, as it was not found in parent", name)
	}

	c[name] = configValue{Value: pv.Value, isSet: true}
	return nil
}

// Merge attempts to inherit any values from `p` that would satisfy required fields.
func (c config) Merge(p config) error {
	for k, v := range c {
		if v.Inherit {
			if verbose {
				fmt.Printf("%q is inherited. Overwriting value with parent's value.\n", k)
			}

			err := c.Inherit(k, p)
			if err != nil {
				return err
			}
		}

		// Inherit the value from the parent, if it's set.
		if v.Required && v.Value == "" {
			if verbose {
				fmt.Printf("%q is required but not given. Inheriting from parent.\n", k)
			}

			err := c.Inherit(k, p)
			if err != nil {
				return fmt.Errorf("Required variable %q, not found locally, or in parent config", k)
			}
		}
	}

	return nil
}

// Environ returns a []string suitable for calling os.Exec with.
func (c config) Environ() []string {
	var kv []string

	return kv
}

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

func run(c config, args []string) {
	if len(args) == 0 {
		// TODO: dump config
		os.Exit(0) // We were just running a check.
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = c.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: fatal: failed to run command: %s\n", os.Args[0], err)
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

func main() {
	var cmd []string

	switch {
	case len(os.Args) < 2:
		usage()
		os.Exit(1)
	case os.Args[1] == "--help":
		usage()
		os.Exit(1)
	}

	var file string

	local := make(config)
	parent := make(config)

	// TODO: Allow local environment to not be read, like in GNU env(1)

	// Read parent, and adjust our assumptions for local + command
	if os.Args[1] == "--stdin" || os.Args[1] == "-" {
		parent.FromEnviron()
		err := parent.Read(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: error: Unable to read parent config from STDIN: %s\n", os.Args[0], err)
			os.Exit(1)
		}

		file = os.Args[2]
		cmd = os.Args[3:]
	} else {
		parent.FromEnviron()
		file = os.Args[1]
		cmd = os.Args[2:]
	}

	// Read local from the `file` we just setup.
	err := local.FromFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: error: Unable to read config from %q: %s\n", os.Args[0], file, err)
		os.Exit(1)
	}

	// Validate that we have what we need.
	err = local.Merge(parent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: error: Environment check failed: %s\n", os.Args[0], err)
		os.Exit(1)
	}

	run(local, cmd)
}
