package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/akamensky/argparse"
)

type Arguments struct {
	printVersion     bool
	verbose          bool
	yesToAll         bool
	zone             int
	installationPath string
	agentcoreToken   string
	agentToken       string
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

	yesToAll := parser.Flag(
		"y",
		"yes",
		&argparse.Options{
			Required: false,
			Help:     "No confirmation prompts",
		},
	)

	installationPath := parser.String(
		"i",
		"installation-path",
		&argparse.Options{
			Required: false,
			Help:     "Path to store the asset Id (not required when an asset Id is provided)",
		},
	)
	agentcoreToken := parser.String(
		"c",
		"agentcore-token",
		&argparse.Options{
			Required: false,
			Help:     "Token for the Agentcore. Must be a container token with `CoreConnect` permissions",
		},
	)
	agentToken := parser.String(
		"a",
		"agent-token",
		&argparse.Options{
			Required: false,
			Help:     "Token for the agents. Must be a container token with  `Read`, `InsertCheckData`, `AssetManagement` and `API` permissions",
		},
	)

	zone := parser.Int(
		"z",
		"zone",
		&argparse.Options{
			Required: false,
			Help:     "Zone Id between 0 and 9 for the agentcore",
			Validate: func(args []string) error {
				if zone, err := strconv.Atoi(args[0]); err == nil {
					if zone < 0 || zone > 9 {
						return errors.New("expecting a value between 0 and 9")
					}
				}
				return nil
			},
			Default: 0,
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
		yesToAll:         *yesToAll,
		zone:             *zone,
		installationPath: *installationPath,
		agentcoreToken:   *agentcoreToken,
		agentToken:       *agentToken,
	}, nil
}

func (args *Arguments) Printf(format string, a ...any) {
	if args.verbose {
		fmt.Printf(format, a...)
	}
}

func (args *Arguments) Println(a ...any) {
	if args.verbose {
		fmt.Println(a...)
	}
}

func (args *Arguments) InstallationPath() {
	if args.installationPath == "" {
		installationPath := "infrasonar"

		if os.Geteuid() == 0 {
			installationPath = "/etc/infrasonar"
		} else if homeDir, err := os.UserHomeDir(); err != nil {
			installationPath = path.Join(homeDir, "infrasonar")
		}
		fmt.Println(installationPath)
	}
}
