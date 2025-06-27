# Performance Tuning Guide: CloudBridge Relay Client v2.0

## Recommended Settings
- Use a fast, local relay server for minimal latency
- Set `relay.timeout` to a reasonable value (e.g., 30s)
- Use TLS session resumption if supported by server
- Tune `rate_limiting` parameters for your workload
- Enable Prometheus metrics for real-time monitoring
- Use the `performance` section in config for runtime optimization

## System Resources
- Ensure enough file descriptors for many tunnels: `ulimit -n 4096`
- Monitor CPU and memory usage under load
- Use `performance.optimization_mode: high_throughput` for max speed, `low_latency` for minimal delay
- Enable `performance.memory_ballast` for better GC performance on large servers

## Buffer Management
- TunnelManager uses BufferManager for efficient data transfer
- Buffer size and pool are optimized for high throughput (default: 4096 bytes, 100 buffers)
- Monitor buffer pool usage via Prometheus metrics
- Error `buffer_pool_exhausted` indicates resource exhaustion

## Logging
- Set `logging.level` to `warn` or `error` in production to reduce overhead
- Use `logging.format: json` for easier parsing

## Tunnel Performance
- Avoid running too many tunnels on a single client instance
- Use dedicated network interfaces for high-throughput tunnels
- Monitor tunnel latency and packet loss
- Per-tunnel statistics are available via Prometheus

## Heartbeat
- Tune heartbeat interval (`heartbeat.interval`) for your reliability needs
- Too frequent heartbeats may increase load; too rare may delay failure detection

## TLS
- Use hardware acceleration for cryptography if available
- Keep CA and client certificates in memory (tmpfs) for faster access

## Monitoring
- Integrate with system monitoring (Prometheus, Grafana, etc.)
- Track connection counts, tunnel stats, heartbeat failures
- Use `/metrics` endpoint for real-time stats (see SECURITY.md for access control)

## Performance Section Example
```yaml
performance:
  enabled: true
  optimization_mode: "high_throughput"  # or "low_latency"
  gc_percent: 100
  memory_ballast: true
```

## Troubleshooting Performance
- Use `--verbose` and system tools (`top`, `htop`, `iftop`) to diagnose bottlenecks
- Check relay server logs for slow responses
- Monitor Prometheus metrics for spikes in buffer usage, errors, or latency 