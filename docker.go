package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const MinDockerVersion = 24
const MinDockerComposeVersion = 2

var reDockerComposeVersion = regexp.MustCompile(`Docker Compose version v(?P<Major>[0-9]+)\.(?P<Minor>[0-9]+)\.(?P<Patch>[0-9]+).*`)
var reDockerVersion = regexp.MustCompile(`Docker version (?P<Major>[0-9]+)\.(?P<Minor>[0-9]+)\.(?P<Patch>[0-9]+).*`)

func dockerComposeVersionCheck(args *Arguments) error {
	out, err := exec.Command("docker", "compose", "version").Output()
	if err != nil {
		return err
	}

	m := reDockerComposeVersion.FindStringSubmatch(string(out))
	if len(m) != 4 {
		if args.ignoreVersionCheck {
			log.Printf("unable to find docker compose version in output: %s", out)
			return nil
		}
		return fmt.Errorf("unable to find docker compose version in output: %s (use --ignore-version-check to ignore this error)", out)
	}

	major, _ := strconv.Atoi(m[1])
	minor, _ := strconv.Atoi(m[2])
	patch, _ := strconv.Atoi(m[3])

	version := fmt.Sprintf("%d.%d.%d", major, minor, patch)

	if major < MinDockerComposeVersion {
		if args.ignoreVersionCheck {
			log.Printf("docker compose version too old: %s (required >= %d.0.0) (ignored)", version, MinDockerVersion)
			return nil
		}
		return fmt.Errorf("docker compose version too old: %s (required >= %d.0.0) (use --ignore-version-check to ignore this error)", version, MinDockerVersion)
	}

	args.Printf("Docker compose version v%s\n", version)
	return nil
}

func dockerVersionCheck(args *Arguments) error {
	out, err := exec.Command("docker", "-v").Output()
	if err != nil {
		return err
	}

	m := reDockerVersion.FindStringSubmatch(string(out))
	if len(m) != 4 {
		if args.ignoreVersionCheck {
			log.Printf("unable to find docker version in output: %s (ignored)", out)
			return nil
		}
		return fmt.Errorf("unable to find docker version in output: %s (use --ignore-version-check to ignore this error)", out)
	}

	major, _ := strconv.Atoi(m[1])
	minor, _ := strconv.Atoi(m[2])
	patch, _ := strconv.Atoi(m[3])

	version := fmt.Sprintf("%d.%d.%d", major, minor, patch)

	if major < MinDockerVersion {
		if args.ignoreVersionCheck {
			log.Printf("docker version too old: %s (required >= %d.0.0) (ignored)", version, MinDockerVersion)
			return nil
		}
		return fmt.Errorf("docker version too old: %s (required >= %d.0.0) (use --ignore-version-check to ignore this error)", version, MinDockerVersion)
	}

	args.Printf("Docker version v%s\n", version)
	return nil
}

func dockerRunningCheck() error {
	out, err := exec.Command("docker", "compose", "ls", "-q").CombinedOutput()
	if err != nil {
		msg := string(out)
		return fmt.Errorf("failed to execute docker compose: %s (%s)", msg, err)
	}
	lines := strings.Split(strings.ReplaceAll(string(out), "\r\n", "\n"), "\n")
	for _, line := range lines {
		if line == "infrasonar" {
			return errors.New("infrasonar appliance already running")
		}
	}
	return nil
}

func dockerRun(cmd *exec.Cmd, args *Arguments) error {
	// Get a pipe to read from standard out
	r, _ := cmd.StdoutPipe()

	// Use the same pipe for standard error
	cmd.Stderr = cmd.Stdout

	// Make a new channel which will be used to ensure we get all output
	done := make(chan bool)

	// Create a scanner which scans r in a line-by-line fashion
	scanner := bufio.NewScanner(r)

	// Use the scanner to scan the output line by line and log it
	// It's running in a goroutine so that it doesn't block
	go func() {

		// Read line by line and process it
		for scanner.Scan() {
			line := scanner.Text()
			args.Printf("%s\n", line)
		}

		// We're all done, unblock the channel
		done <- true

	}()

	// Start the command and check for errors
	err := cmd.Start()
	if err != nil {
		return err
	}

	// Wait for all output to be processed
	<-done

	// Wait for the command to finish
	err = cmd.Wait()
	return err
}

func dockerStart(args *Arguments) error {
	cmd := exec.Command("docker", "compose", "--progress", "plain", "pull")
	cmd.Dir = args.installationPath
	args.Printf("Pulling images...\n")
	if err := dockerRun(cmd, args); err != nil {
		return err
	}

	cmd = exec.Command("docker", "compose", "--progress", "plain", "up", "-d")
	cmd.Dir = args.installationPath
	args.Printf("Starting containers...\n")
	if err := dockerRun(cmd, args); err != nil {
		return err
	}

	return nil
}
