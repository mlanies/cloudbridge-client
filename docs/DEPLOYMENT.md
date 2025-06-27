# Deployment Guide: CloudBridge Relay Client v2.0

## Prerequisites
- Go 1.21+
- Access to relay server (TLS 1.3 required)
- Valid JWT token with tenant_id or Keycloak credentials
- For metrics: Prometheus/Grafana (optional)

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
export CLOUDBRIDGE_METRICS_ENABLED="true"
export CLOUDBRIDGE_PERFORMANCE_ENABLED="true"
```

## Configuration Example (v2.0)
```yaml
relay:
  host: "edge.2gc.ru"
  port: 8080
  timeout: "30s"
  tls:
    enabled: true
    min_version: "1.3"

auth:
  type: "jwt"
  secret: "your-jwt-secret-here"

metrics:
  enabled: true
  prometheus_port: 9090
  tenant_metrics: true
  buffer_metrics: true

performance:
  enabled: true
  optimization_mode: "high_throughput"
  gc_percent: 100
  memory_ballast: true
```

## Systemd Service Example
Create `/etc/systemd/system/cloudbridge-client.service`:
```ini
[Unit]
Description=CloudBridge Relay Client v2.0
After=network.target

[Service]
ExecStart=/path/to/cloudbridge-client --config /path/to/config.yaml --token "your-jwt-token"
Restart=on-failure
User=ubuntu
# For metrics endpoint security
NoNewPrivileges=true
PrivateNetwork=false

[Install]
WantedBy=multi-user.target
```

## Prometheus Integration
If metrics are enabled, the client exposes metrics on `http://localhost:9090/metrics`:
```bash
# Test metrics endpoint
curl http://localhost:9090/metrics

# Add to prometheus.yml
scrape_configs:
  - job_name: 'cloudbridge-client'
    static_configs:
      - targets: ['localhost:9090']
```

## Performance Tuning
- Use `performance.optimization_mode: "high_throughput"` for maximum speed
- Use `performance.optimization_mode: "low_latency"` for minimal delay
- Enable `performance.memory_ballast` on large servers
- Monitor buffer usage via Prometheus metrics

## Multi-Tenancy Setup
- Ensure JWT tokens include `tenant_id` claim
- Configure tenant-specific limits on relay server
- Monitor per-tenant metrics via Prometheus

## Updating
- Pull latest changes: `git pull`
- Rebuild: `go build -o cloudbridge-client ./cmd/cloudbridge-client`
- Restart service: `sudo systemctl restart cloudbridge-client`

## Logs
- By default, logs are printed to stdout. Configure via `config.yaml`.
- All operations include tenant_id for multi-tenancy environments.

## Troubleshooting
- See `docs/TROUBLESHOOTING.md` for common issues.
- Check Prometheus metrics for performance bottlenecks.
- Verify tenant_id is present in JWT tokens. 