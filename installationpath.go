package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ensureInstallationPath(args *Arguments) error {
	if args.installationPath == "" {
		mydir, err := os.Getwd()
		if err != nil {
			return err
		}
		installationPath := filepath.Join(mydir, "infrasonar")

		if os.Geteuid() == 0 {
			installationPath = "/etc/infrasonar"
		} else if homeDir, err := os.UserHomeDir(); err == nil {
			installationPath = filepath.Join(homeDir, "infrasonar")
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
		args.installationPath = filepath.Join(args.installationPath, "infrasonar")
	}
	configPath := filepath.Join(args.installationPath, "data", "config")

	_, err := os.Stat(args.installationPath)

	if err == nil && !args.yesToAll {
		fmt.Printf("A file or directory already exists at '%s'.\n\nOverwrite and continue? (yes/no)\n", args.installationPath)
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
