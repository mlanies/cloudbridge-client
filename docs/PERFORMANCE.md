# Performance Tuning Guide: CloudBridge Relay Client

## Recommended Settings
- Use a fast, local relay server for minimal latency
- Set `relay.timeout` to a reasonable value (e.g., 30s)
- Use TLS session resumption if supported by server
- Tune `rate_limiting` parameters for your workload

## System Resources
- Ensure enough file descriptors for many tunnels: `ulimit -n 4096`
- Monitor CPU and memory usage under load

## Logging
- Set `logging.level` to `warn` or `error` in production to reduce overhead
- Use `logging.format: json` for easier parsing

## Tunnel Performance
- Avoid running too many tunnels on a single client instance
- Use dedicated network interfaces for high-throughput tunnels
- Monitor tunnel latency and packet loss

## Heartbeat
- Tune heartbeat interval (`heartbeat.interval`) for your reliability needs
- Too frequent heartbeats may increase load; too rare may delay failure detection

## TLS
- Use hardware acceleration for cryptography if available
- Keep CA and client certificates in memory (tmpfs) for faster access

## Monitoring
- Integrate with system monitoring (Prometheus, Grafana, etc.)
- Track connection counts, tunnel stats, heartbeat failures

## Troubleshooting Performance
- Use `--verbose` and system tools (`top`, `htop`, `iftop`) to diagnose bottlenecks
- Check relay server logs for slow responses 