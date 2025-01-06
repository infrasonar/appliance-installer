package main

import (
	"fmt"
	"os"
)

func exitOnErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	// Read the arguments
	args, err := parseArgs()
	exitOnErr(err)

	if args.printVersion {
		fmt.Printf("InfraSonar appliance installer v%s\n", Version)
		os.Exit(0)
	}

	if os.Geteuid() > 0 && !args.yesToAll {
		// UID < 0 is Windows, UID = 0 is root on Linux and UID > 0 is a user
		fmt.Println("The InfraSonar Appliance is typically installed as root. Are you sure you want to continue as a normal user? (yes/no)")
		if !askForConfirmation() {
			exitOnErr(ErrUserCanceled)
		}
	}

	// Check if docker is installed and no appliance is running
	exitOnErr(dockerVersionCheck(args))
	exitOnErr(dockerComposeVersionCheck(args))
	exitOnErr(dockerRunningCheck())

	// Ensure we have an installation path
	exitOnErr(ensureInstallationPath(args))

	// Ask for required token
	args.EnsureAgentcoreToken()
	args.EnsureAgentToken()

	if !args.yesToAll {
		fmt.Printf(`
--------------------------------------------------------------------------
  The appliance will be installed in '%s' for zone %d
--------------------------------------------------------------------------

Do you want to continue? (yes/no)
`, args.installationPath, args.zone)
		if !askForConfirmation() {
			exitOnErr(ErrUserCanceled)
		}
		fmt.Println("Please be patient, this may take a while...")
	}

	// Installation
	exitOnErr(install(args))

	// Start InfraSonar appliance
	exitOnErr(dockerStart(args))

	if !args.yesToAll {
		fmt.Printf("done\n")
	} else {
		args.Printf("done\n")
	}
}
