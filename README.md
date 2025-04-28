[![CI](https://github.com/infrasonar/appliance-installer/workflows/CI/badge.svg)](https://github.com/infrasonar/appliance-installer/actions)
[![Release Version](https://img.shields.io/github/release/infrasonar/appliance-installer)](https://github.com/infrasonar/appliance-installer/releases)

# InfraSonar Appliance Installer

This tool facilitates the deployment of the InfraSonar appliance.

**Prerequisites:** Docker version 24 or higher is required.

## Installation

**1. Download the latest installer:**

- [Linux (amd64)](https://github.com/infrasonar/appliance-installer/releases/download/v1.0.8/appliance-installer-linux-amd64-1.0.8.tar.gz)
- [Darwin (amd64)](https://github.com/infrasonar/appliance-installer/releases/download/v1.0.8/appliance-installer-darwin-amd64-1.0.8.tar.gz)
- [Windows (amd64)](https://github.com/infrasonar/appliance-installer/releases/download/v1.0.8/appliance-installer-windows-amd64-1.0.8.zip)
- [Solaris (amd64)](https://github.com/infrasonar/appliance-installer/releases/download/v1.0.8/appliance-installer-solaris-amd64-1.0.8.tar.gz)

> If your platform is not listed above, refer to the [build from source](#build-from-source) section for instructions.

**2. Extract the contents of the archive using a tool like `tar`. Here's an example for Linux (amd64):**
```bash
tar -xzvf appliance-installer-linux-amd64-1.0.8.tar.gz
```

**3. Run the installer:**

```bash
./appliance-installer --verbose
```

## Help information

```
usage: installer [-h|--help] [-i|--installation-path "<value>"]
                 [-c|--agentcore-token "<value>"] [-a|--agent-token "<value>"]
                 [-z|--zone <integer>] [-d|--use-development] [-v|--verbose]
                 [-y|--yes] [--version]

                 Installation of the InfraSonar appliance

Arguments:

  -h  --help                  Print help information
  -i  --installation-path     Installation path for the docker compose and
                              configuration files
  -c  --agentcore-token       Token for the Agentcore. Must be a container
                              token with `CoreConnect` permissions
  -a  --agent-token           Token for the agents. Must be a container token
                              with `Read`, `InsertCheckData`, `AssetManagement`
                              and `API` permissions
  -z  --zone                  Zone Id between 0 and 9. Default: 0
  -r  --remote-access         Enable the option for remote access on this
                              appliance for users with CoreConnect permissions
                              (only applicable with -y/--yes)
  -d  --use-development       Use the InfraSonar development environment
  -v  --verbose               Enable verbose output
  -y  --yes                   No confirmation prompts
      --ignore-version-check  Ignore Docker and Docker Compose version check
      --version               Print version information and quit
```

## Build from source

Make sure [Go](https://go.dev/doc/install) is installed.

**1. Clone the repository**
```bash
git clone https://github.com/infrasonar/appliance-installer.git
```

**2. Open the cloned directory**
```bash
cd appliance-installer
```

**3. Build the appliance installer**
```bash
CGO_ENABLED=0 go build -o appliance-installer
```
