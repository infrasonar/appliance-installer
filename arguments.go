package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/fatih/color"
)

var reToken = regexp.MustCompile(`^[0-9a-f]{32}$`)
var tokenValidation = func(args []string) error {
	if !reToken.MatchString(args[0]) {
		return errors.New("invalid token")
	}
	return nil
}

func askToken() string {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		response = ""
	}
	response = strings.TrimSpace(response)
	if reToken.MatchString(response) {
		return response
	} else {
		fmt.Println("Invalid token, please enter a correct token")
		return askToken()
	}
}

type Arguments struct {
	printVersion       bool
	verbose            bool
	yesToAll           bool
	noRemoteAccess     bool
	useDevelopment     bool
	ignoreVersionCheck bool
	zone               int
	installationPath   string
	agentcoreToken     string
	agentToken         string
}

func parseArgs() (*Arguments, error) {
	parser := argparse.NewParser("installer", "Installation of the InfraSonar appliance")
	installationPath := parser.String(
		"i",
		"installation-path",
		&argparse.Options{
			Required: false,
			Help:     "Installation path for the docker compose and configuration files",
		},
	)
	agentcoreToken := parser.String(
		"c",
		"agentcore-token",
		&argparse.Options{
			Required: false,
			Validate: tokenValidation,
			Help:     "Token for the Agentcore. Must be a container token with `CoreConnect` permissions",
		},
	)
	agentToken := parser.String(
		"a",
		"agent-token",
		&argparse.Options{
			Required: false,
			Validate: tokenValidation,
			Help:     "Token for the agents. Must be a container token with `Read`, `InsertCheckData`, `AssetManagement` and `API` permissions",
		},
	)
	zone := parser.Int(
		"z",
		"zone",
		&argparse.Options{
			Required: false,
			Help:     "Zone Id between 0 and 9",
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
	noRemoteAccess := parser.Flag(
		"r",
		"no-remote-access",
		&argparse.Options{
			Required: false,
			Help:     "Disable remote access on this appliance",
		},
	)
	useDevelopment := parser.Flag(
		"d",
		"use-development",
		&argparse.Options{
			Required: false,
			Help:     "Use the InfraSonar development environment",
		},
	)
	verbose := parser.Flag(
		"v",
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
	ignoreVersionCheck := parser.Flag(
		"",
		"ignore-version-check",
		&argparse.Options{
			Required: false,
			Help:     "Ignore Docker and Docker Compose version check",
		},
	)
	printVersion := parser.Flag(
		"",
		"version",
		&argparse.Options{
			Required: false,
			Help:     "Print version information and quit",
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
		printVersion:       *printVersion,
		verbose:            *verbose,
		yesToAll:           *yesToAll,
		noRemoteAccess:     *noRemoteAccess,
		useDevelopment:     *useDevelopment,
		ignoreVersionCheck: *ignoreVersionCheck,
		zone:               *zone,
		installationPath:   *installationPath,
		agentcoreToken:     *agentcoreToken,
		agentToken:         *agentToken,
	}, nil
}

func (args *Arguments) Printf(format string, a ...any) {
	if args.verbose {
		color.Yellow(format, a...)
	}
}

func (args *Arguments) EnsureAgentcoreToken() {
	if args.agentcoreToken == "" {
		fmt.Println("Please provide a token for the Agentcore (container token with `CoreConnect` permissions):")
		args.agentcoreToken = askToken()
	}
}

func (args *Arguments) EnsureAgentToken() {
	if args.agentToken == "" {
		fmt.Println("Please provide a token for the agents (container token with `Read`, `InsertCheckData`, `AssetManagement` and `API` permissions):")
		args.agentToken = askToken()
	}
}

func (args *Arguments) EnsureRemoteAccess() {
	if !args.noRemoteAccess && !args.yesToAll {
		fmt.Println("Do you want to enable the option for users with CoreConnect permissions to start remote access? (yes/no)")
		if askForConfirmation() {
			args.noRemoteAccess = true
		} else {
			args.noRemoteAccess = false
		}
	}
}
