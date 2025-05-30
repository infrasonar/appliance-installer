package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const templateEnv = `AGENTCORE_TOKEN = <AGENTCORE_TOKEN>
AGENT_TOKEN = <AGENT_TOKEN>
AGENTCORE_ZONE_ID = <AGENTCORE_ZONE_ID>
SOCAT_TARGET_ADDR = ""
`

const templateCompose = `## InfraSonar docker-compose.yml file
##
## !! This file is managed by InfraSonar !!

services:
  agentcore:
    environment:
      AGENTCORE_ZONE: ${AGENTCORE_ZONE_ID}
      HUB_HOST: <HUB_ADDRESS>
      LOG_COLORIZED: 1
      LOG_LEVEL: info
      TOKEN: ${AGENTCORE_TOKEN}
    image: ghcr.io/infrasonar/agentcore
    logging: &id001
      options:
        max-size: 5m
    network_mode: host
    restart: always
    volumes: &id002
    - <DATA_PATH>:/data/
  rapp:
    environment:
      ALLOW_REMOTE_ACCESS: <ALLOW_REMOTE_ACCESS>
      COMPOSE_FILE: /etc/infrasonar/docker-compose.yml
      CONFIG_FILE: /etc/infrasonar/data/config/infrasonar.yaml
      ENV_FILE: /etc/infrasonar/.env
      DATA_PATH: <DATA_PATH>
      USE_DEVELOPMENT: <USE_DEVELOPMENT>
    image: ghcr.io/infrasonar/rapp
    logging: *id001
    network_mode: host
    restart: always
    volumes:
    - <INSTALLATION_PATH>:/etc/infrasonar
    - /var/run/docker.sock:/var/run/docker.sock
x-infrasonar-template:
  logging: *id001
  network_mode: host
  restart: always
  volumes: *id002
`

const configContent = `## WARNING: InfraSonar will make 'password' and 'secret' values unreadable but
## this must not be regarded as true encryption as the encryption key is
## publicly available.
##
## Example configuration for 'myprobe' collector:
##
##  myprobe:
##    config:
##      username: alice
##      password: "secret password"
##    assets:
##    - id: [12345, 34567]
##      config:
##        username: bob
##        password: "my secret"
##
## !! This file is managed by InfraSonar !!
##
## It's okay to add custom probe configuration for when you want to
## specify the "_use" value for assets. The appliance manager will not
## overwrite these custom probe configurations. You can also add additional
## assets configurations for managed probes.
`

func installEnv(args *Arguments) error {
	content := strings.Replace(templateEnv, "<AGENTCORE_TOKEN>", args.agentcoreToken, 1)
	content = strings.Replace(content, "<AGENT_TOKEN>", args.agentToken, 1)
	content = strings.Replace(content, "<AGENTCORE_ZONE_ID>", fmt.Sprint(args.zone), 1)
	fn := filepath.Join(args.installationPath, ".env")
	fp, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer fp.Close()
	_, err = fp.WriteString(content)
	return err
}

func installCompose(args *Arguments) error {
	hub_address := "hub.infrasonar.com"
	use_development := "0"
	if args.useDevelopment {
		hub_address = "devhub.infrasonar.com"
		use_development = "1"
	}
	allow_remote_access := "0"
	if args.allowRemoteAccess {
		allow_remote_access = "1"
	}
	content := strings.Replace(templateCompose, "<HUB_ADDRESS>", hub_address, 1)
	content = strings.ReplaceAll(content, "<INSTALLATION_PATH>", args.installationPath)
	content = strings.Replace(content, "<ALLOW_REMOTE_ACCESS>", allow_remote_access, 1)
	content = strings.Replace(content, "<USE_DEVELOPMENT>", use_development, 1)
	content = strings.ReplaceAll(content, "<DATA_PATH>", fmt.Sprintf(".%cdata", os.PathSeparator))

	fn := filepath.Join(args.installationPath, "docker-compose.yml")
	fp, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer fp.Close()
	_, err = fp.WriteString(content)
	return err
}

func installConfig(args *Arguments) error {
	fn := filepath.Join(args.installationPath, "data", "config", "infrasonar.yaml")

	fp, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer fp.Close()
	_, err = fp.WriteString(configContent)
	return err
}

func install(args *Arguments) error {
	if err := installEnv(args); err != nil {
		return err
	}
	if err := installCompose(args); err != nil {
		return err
	}
	if err := installConfig(args); err != nil {
		return err
	}
	return nil
}
