package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// Read the arguments
	args, err := parseArgs()
	if err != nil {
		log.Fatal(err)
	}

	if args.printVersion {
		fmt.Printf("InfraSonar appliance installer v%s\n", Version)
		os.Exit(0)
	}

	if err := dockerVersionCheck(args); err != nil {
		log.Fatal(err)
	}
	if err := dockerComposeVersionCheck(args); err != nil {
		log.Fatal(err)
	}

}
