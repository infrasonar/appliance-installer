[![CI](https://github.com/infrasonar/appliance-installer/workflows/CI/badge.svg)](https://github.com/infrasonar/appliance-installer/actions)
[![Release Version](https://img.shields.io/github/release/infrasonar/appliance-installer)](https://github.com/infrasonar/appliance-installer/releases)

# InfraSonar Appliance Installer

## Help information

```
usage: installer [-h|--help] [-i|--installation-path "<value>"]
                 [-c|--agentcore-token "<value>"] [-a|--agent-token "<value>"]
                 [-z|--zone <integer>] [-d|--use-development] [-v|--verbose]
                 [-y|--yes] [--version]

                 Installation of the InfraSonar appliance

Arguments:

  -h  --help               Print help information
  -i  --installation-path  Installation path for the docker compose and
                           configuration files
  -c  --agentcore-token    Token for the Agentcore. Must be a container token
                           with `CoreConnect` permissions
  -a  --agent-token        Token for the agents. Must be a container token with
                           `Read`, `InsertCheckData`, `AssetManagement` and
                           `API` permissions
  -z  --zone               Zone Id between 0 and 9. Default: 0
  -d  --use-development    Use the InfraSonar development environment
  -v  --verbose            Enable verbose output
  -y  --yes                No confirmation prompts
      --version            Print version information and quit
```

## Downloads




## Build
```
CGO_ENABLED=0 go build -o appliance-installer
```
