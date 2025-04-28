package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func exitOnErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	// Read the arguments
	args, err := parseArgs()
	exitOnErr(err)

	// Check if we just need to print the version and exit
	if args.printVersion {
		fmt.Printf("InfraSonar appliance installer v%s\n", Version)
		os.Exit(0)
	}

	// Check if docker is installed and no appliance is running
	exitOnErr(dockerVersionCheck(args))
	exitOnErr(dockerComposeVersionCheck(args))
	exitOnErr(dockerRunningCheck())

	// Ensure we have an installation path
	exitOnErr(ensureInstallationPath(args))

	// Ask for required token and remote access
	args.EnsureRemoteAccess()
	args.EnsureAgentcoreToken()
	args.EnsureAgentToken()

	if !args.yesToAll {
		content := fmt.Sprintf("The appliance for zone %d will be deployed in the '%s' directory", args.zone, args.installationPath)
		divider := strings.Repeat("#", len(content)+4)

		color.Yellow("\n%s\n\n  %s\n\n%s\n\n", divider, content, divider)
		fmt.Println("Do you want to continue? (yes/no)")
		if !askForConfirmation() {
			exitOnErr(ErrUserCanceled)
		}
		fmt.Println("Please be patient, this may take a while...")
	}

	// Installation
	exitOnErr(install(args))

	// Start InfraSonar appliance
	exitOnErr(dockerStart(args))

	website := "https://app.infrasonar.com"
	if args.useDevelopment {
		website = "https://devapp.infrasonar.com"
	}

	if !args.yesToAll {
		fmt.Printf("Done\n")
		fmt.Printf("Open your container on %s and manage the appliance via the 'Agentcores' menu\n", website)
	} else {
		args.Printf("done\n")
	}
}
