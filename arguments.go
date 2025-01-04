package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

type Arguments struct {
	printVersion     bool
	verbose          bool
	installationPath *string
	agentcoreToken   *string
	agentToken       *string
}

func parseArgs() (*Arguments, error) {
	parser := argparse.NewParser("installer", "Installation of the InfraSonar appliance")
	printVersion := parser.Flag(
		"v",
		"version",
		&argparse.Options{
			Required: false,
			Help:     "Print version information and quit",
		},
	)

	verbose := parser.Flag(
		"",
		"verbose",
		&argparse.Options{
			Required: false,
			Help:     "Enable verbose output",
		},
	)

	installationPath := parser.String(
		"i",
		"installation-path",
		&argparse.Options{
			Required: false,
			Help:     "Path to store the asset Id (not required when an asset Id is provided)",
			Default:  "",
		},
	)
	agentcoreToken := parser.String(
		"c",
		"agentcore-token",
		&argparse.Options{
			Required: false,
			Help:     "Token for the Agentcore. Must be a container token with `CoreConnect` permissions",
			Default:  "",
		},
	)
	agentToken := parser.String(
		"a",
		"agent-token",
		&argparse.Options{
			Required: false,
			Help:     "Token for the agents. Must be a container token with  `Read`, `InsertCheckData`, `AssetManagement` and `API` permissions",
			Default:  "",
		},
	)

	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		return nil, err
	}

	return &Arguments{
		printVersion:     *printVersion,
		verbose:          *verbose,
		installationPath: installationPath,
		agentcoreToken:   agentcoreToken,
		agentToken:       agentToken,
	}, nil
}

func (args *Arguments) Printf(format string, a ...any) {
	if args.verbose {
		fmt.Printf(format, a...)
	}
}
