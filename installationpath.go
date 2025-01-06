package main

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func ensureInstallationPath(args *Arguments) error {
	if args.installationPath == "" {
		mydir, err := os.Getwd()
		if err != nil {
			return err
		}
		installationPath := path.Join(mydir, "infrasonar")

		if os.Geteuid() == 0 {
			installationPath = "/etc/infrasonar"
		} else if homeDir, err := os.UserHomeDir(); err == nil {
			installationPath = path.Join(homeDir, "infrasonar")
		}
		fmt.Printf("Installation Path (enter path or press Enter for default: %s)\n", installationPath)
		var response string
		_, err = fmt.Scanln(&response)
		if err == nil {
			response = strings.TrimSpace(response)
			if response == "" {
				args.installationPath = installationPath
			} else {
				args.installationPath = response
			}
		} else {
			args.installationPath = installationPath
		}
	}

	if !strings.HasSuffix(args.installationPath, "infrasonar") {
		args.installationPath = path.Join(args.installationPath, "infrasonar")
	}
	configPath := path.Join(args.installationPath, "data", "config")

	_, err := os.Stat(args.installationPath)

	if err == nil && !args.yesToAll {
		fmt.Printf("A file or directory already exists at '%s'. Overwrite and continue? (yes/no)\n", args.installationPath)
		if !askForConfirmation() {
			return ErrUserCanceled
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	err = os.MkdirAll(configPath, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
