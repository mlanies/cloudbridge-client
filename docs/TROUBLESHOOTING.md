# Troubleshooting Guide: CloudBridge Relay Client v2.0

## Common Issues & Solutions

### 1. TLS Handshake Fails
- **Check**: Is the relay server running and reachable?
- **Check**: Are you using TLS 1.3 and correct cipher suites?
- **Check**: Is the CA certificate path correct in config?
- **Check**: Is the server certificate valid and not expired?

### 2. Authentication Fails (`invalid_token`)
- **Check**: Is your JWT token valid (not expired, correct signature)?
- **Check**: Does the JWT secret match the relay server?
- **Check**: For Keycloak, are realm, client_id, and JWKS endpoint correct?
- **Check**: Does the token have a `sub` claim?
- **Check**: For multi-tenancy, does the token have a `tenant_id` claim?

### 3. Multi-Tenancy Issues
- **Error**: `tenant_not_found` - Check if tenant_id is present in JWT token
- **Error**: `tenant_limit_exceeded` - Tenant has reached tunnel/connection limits
- **Error**: `ip_not_allowed` - Client IP not in allowed range for tenant
- **Check**: Verify tenant_id is correctly set in JWT token
- **Check**: Confirm tenant exists and is active on relay server

### 4. Rate Limit Exceeded
- **Check**: Are you sending too many requests per second?
- **Check**: For multi-tenancy, check per-tenant limits
- **Solution**: Wait for backoff and retry, or increase limits in config if you control the server.

### 5. Tunnel Creation Fails
- **Check**: Are local and remote ports valid and not in use?
- **Check**: Is the remote host reachable from the relay server?
- **Check**: Is the tunnel_id unique?
- **Check**: For multi-tenancy, verify tenant_id is included in tunnel_info message

### 6. Performance Issues
- **Error**: `buffer_pool_exhausted` - Too many concurrent connections
- **Error**: `data_transfer_failed` - Network or buffer issues
- **Check**: Monitor Prometheus metrics for buffer usage
- **Solution**: Increase buffer pool size or reduce concurrent connections
- **Check**: Use performance.optimization_mode for your workload

### 7. Metrics Issues
- **Check**: Is metrics.enabled set to true in config?
- **Check**: Is port 9090 available for Prometheus endpoint?
- **Check**: Can you access http://localhost:9090/metrics?
- **Check**: Are metrics being collected (check for cloudbridge_* metrics)?

### 8. Heartbeat Fails
- **Check**: Is the connection to the relay server still alive?
- **Check**: Is there network latency or firewall issues?

### 9. Unknown Message Type
- **Check**: Are you using a compatible client and server version?
- **Check**: Is the protocol version in hello message correct?

## Debugging Tips
- Run with `--verbose` to enable detailed logging.
- Check logs for error codes and messages.
- Use `openssl s_client` to debug TLS connections.
- Validate JWT tokens with [jwt.io](https://jwt.io/).
- Check relay server logs for more details.
- Monitor Prometheus metrics for performance insights.
- Check tenant_id in JWT tokens and tunnel messages.

## Performance Debugging
- Use `curl http://localhost:9090/metrics` to check buffer usage
- Monitor `cloudbridge_buffer_pool_usage` metric
- Check `cloudbridge_bytes_transferred_total` for throughput
- Verify `cloudbridge_active_connections` for connection count

## Getting Help
- Review the README and docs/README.md for configuration and usage.
- Open an issue on GitHub with logs and config details (do not include secrets).
- For security issues, contact the security contact in the main README.
- Include Prometheus metrics in bug reports if relevant. 