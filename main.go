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
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [OPTIONS]... JSONfile [COMMAND [ARG]...]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, `
Ensure that the environment meets JSONfile requirements and run COMMAND.

  -d, --display-docs  displays documentation for variables in spec
  -h, --help     display this help and exit
  -i, --ignore-environment  start with an empty environment
  -,  --stdin    insert read JSON key-value pairs into environment
  -v, --validate-json  validates envjson file

If no COMMAND, print the resulting environment as JSON.
`)
}

type config struct {
	displayHelp  bool
	displayEnv   bool
	displayDocs  bool	
	validateJSON bool
	fromStdin    bool
	cleanEnv     bool
	prog         string
	specFile     string
	cmd          []string
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
		case args[i] == "--help" || args[i] == "-h":
			c.displayHelp = true
			return
		case args[i] == "--stdin" || args[i] == "-":
			c.fromStdin = true
		case args[i] == "-i" || args[i] == "--ignore-environment":
			c.cleanEnv = true
		case args[i] == "-v" || args[i] == "--validate-json":
			c.validateJSON = true
		case args[i] == "-d" || args[i] == "--display-docs":
			c.displayDocs = true
		case strings.HasPrefix(args[i], "--") || strings.HasPrefix(args[i], "-"):
			err = fmt.Errorf("invalid argument %q", args[i])
			return
		default:
			c.displayEnv = true
			break loop
		}
	}

	if i < len(args) {
		c.specFile = args[i] // save off the spec file for loading later.

		i++
		if i < len(args) {
			c.cmd = args[i:]
			c.displayEnv = false
		}
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

	if config.validateJSON {
		ok, err := validateJSON(config.specFile)
		if !ok {
			fmt.Fprintf(os.Stderr,
				"%s: error: Unable to validate %q: %s\n",
				config.prog, config.specFile, err)
			os.Exit(1)
		}
		os.Exit(0)
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
	if config.displayDocs {
		err = local.DumpDocs(os.Stdout)
		if err != nil {
			fmt.Fprintf(os.Stderr,
				"%s: error: Unable to display docs: %s\n",
				config.prog, err)
			os.Exit(1)
		}
		return
	} else if config.displayEnv {
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
