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

	if os.Geteuid() > 0 {
		
	}

	exitOnErr(dockerVersionCheck(args))
	exitOnErr(dockerComposeVersionCheck(args))

	args.Println("done")
}
