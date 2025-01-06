# InfraSonar Appliance Installer



## Build
```
CGO_ENABLED=0 go build -o appliance-installer
```


GOOS=solaris GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -o acsls-agent.solaris-amd64