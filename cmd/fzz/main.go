package main

import (
	"errors"
	"fmt"
	"github.com/foojank/fzz"
	"github.com/foojank/fzz/services/builder"
	"github.com/foojank/fzz/services/executor"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

var version string

func Version() string {
	return strings.TrimSpace(version)
}

var app = cli.App{
	Name:           "fzz",
	Usage:          "Build or run fizzpit files",
	Args:           true,
	Version:        Version(),
	DefaultCommand: "",
	Commands: cli.Commands{
		{
			Name:   "build",
			Usage:  "Build fizzpit from a source directory",
			Args:   true,
			Action: build,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "o",
					Value: "o.fzz",
					Usage: "output file",
				},
			},
			HideHelpCommand: true,
		},
		{
			Name:   "run",
			Usage:  "Run fizzpit file",
			Args:   true,
			Action: run,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "f",
					Required: true,
					Usage:    "input file",
				},
			},
			HideHelpCommand: true,
		},
	},
}

func build(c *cli.Context) error {
	args := c.Args()
	if !args.Present() {
		return errors.New("no argument specified")
	}

	src := args.First()
	out := c.String("o")
	err := fzz.Build(c.Context, src, builder.Arguments{
		Output: out,
	})
	if err != nil {
		return err
	}

	return nil
}

func run(c *cli.Context) error {
	args := c.Args()
	if !args.Present() {
		return errors.New("no argument specified")
	}

	file := c.String("f")
	b, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	cmd := args.Get(0)
	cmdArgs := make([]string, 0, args.Len()-1)
	for i := 1; i < args.Len(); i++ {
		cmdArgs = append(cmdArgs, args.Get(i))
	}

	err = fzz.Exec(c.Context, b, executor.Arguments{
		Command:      cmd,
		Args:         cmdArgs,
		Env:          nil,
		Unrestricted: true, // TODO: disable by default, configurable!
		Imports: []fzz.Imports{
			fzz.Stdlib,
			fzz.Unrestricted, // TODO: disable by default, configurable!
			fzz.Syscall,      // TODO: disable by default, configurable!
			fzz.Unsafe,       // TODO: disable by default, configurable!
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
