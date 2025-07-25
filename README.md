# harbor_honeypot

This is WEB honeypot implementing harbor login UI. Static usage from `https://github.com/goharbor`

### Build

Docker:
```bash
docker built -t harbor_honeypot -f docker/Dockerfile
```

From src:
```bash
go build -o harbor_honeypot
```

### Environment
1. APP_PORT - http port. Default is `8080`.
2. APP_LOG_FILE_PATH - path to log file. Default is `/var/log/honeypot/harbor_honeypot.json`.