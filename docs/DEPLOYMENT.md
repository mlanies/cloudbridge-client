# Deployment Guide: CloudBridge Relay Client

## Prerequisites
- Go 1.20+
- Access to relay server (TLS 1.3 required)
- Valid JWT token or Keycloak credentials

## Building from Source
```bash
git clone https://github.com/2gc-dev/cloudbridge-client.git
cd cloudbridge-client
go build -o cloudbridge-client ./cmd/cloudbridge-client
```

## Pre-built Binaries
- Download from the [releases page](https://github.com/2gc-dev/cloudbridge-client/releases)
- Make executable: `chmod +x cloudbridge-client`

## Running the Client
```bash
./cloudbridge-client --token "your-jwt-token"
```

## Using a Configuration File
```bash
./cloudbridge-client --config config.yaml --token "your-jwt-token"
```

## Environment Variables
All config options can be set via `CLOUDBRIDGE_` prefix, e.g.:
```bash
export CLOUDBRIDGE_RELAY_HOST="relay.example.com"
export CLOUDBRIDGE_AUTH_SECRET="your-jwt-secret"
```

## Systemd Service Example
Create `/etc/systemd/system/cloudbridge-client.service`:
```ini
[Unit]
Description=CloudBridge Relay Client
After=network.target

[Service]
ExecStart=/path/to/cloudbridge-client --config /path/to/config.yaml --token "your-jwt-token"
Restart=on-failure
User=ubuntu

[Install]
WantedBy=multi-user.target
```

## Updating
- Pull latest changes: `git pull`
- Rebuild: `go build -o cloudbridge-client ./cmd/cloudbridge-client`

## Logs
- By default, logs are printed to stdout. Configure via `config.yaml`.

## Troubleshooting
- See `docs/TROUBLESHOOTING.md` for common issues. 