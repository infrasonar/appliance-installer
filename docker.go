package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

const MinDockerVersion = 24
const MinDockerComposeVersion = 2

func dockerComposeVersionCheck(args *Arguments) error {
	reDockerComposeVersion := regexp.MustCompile(`Docker Compose version v(?P<Major>[0-9]+)\.(?P<Minor>[0-9]+)\.(?P<Patch>[0-9]+).*`)

	out, err := exec.Command("docker", "compose", "version").Output()
	if err != nil {
		return err
	}

	m := reDockerComposeVersion.FindStringSubmatch(string(out))
	if len(m) != 4 {
		return fmt.Errorf("unable to find docker compose version in output: %s", out)
	}

	major, _ := strconv.Atoi(m[1])
	minor, _ := strconv.Atoi(m[2])
	patch, _ := strconv.Atoi(m[3])

	version := fmt.Sprintf("%d.%d.%d", major, minor, patch)

	if major < MinDockerComposeVersion {
		return fmt.Errorf("docker compose version too old: %s (required >= %d.0.0)", version, MinDockerVersion)
	}

	args.Printf("docker compose version v%s\n", version)
	return nil
}

func dockerVersionCheck(args *Arguments) error {
	reDockerVersion := regexp.MustCompile(`Docker version (?P<Major>[0-9]+)\.(?P<Minor>[0-9]+)\.(?P<Patch>[0-9]+).*`)

	out, err := exec.Command("docker", "-v").Output()
	if err != nil {
		return err
	}

	m := reDockerVersion.FindStringSubmatch(string(out))
	if len(m) != 4 {
		return fmt.Errorf("unable to find docker version in output: %s", out)
	}

	major, _ := strconv.Atoi(m[1])
	minor, _ := strconv.Atoi(m[2])
	patch, _ := strconv.Atoi(m[3])

	version := fmt.Sprintf("%d.%d.%d", major, minor, patch)

	if major < MinDockerVersion {
		return fmt.Errorf("docker version too old: %s (required >= %d.0.0)", version, MinDockerVersion)
	}

	args.Printf("docker version v%s\n", version)
	return nil
}
