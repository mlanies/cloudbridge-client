# Troubleshooting Guide: CloudBridge Relay Client

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

### 3. Rate Limit Exceeded
- **Check**: Are you sending too many requests per second?
- **Solution**: Wait for backoff and retry, or increase limits in config if you control the server.

### 4. Tunnel Creation Fails
- **Check**: Are local and remote ports valid and not in use?
- **Check**: Is the remote host reachable from the relay server?
- **Check**: Is the tunnel_id unique?

### 5. Heartbeat Fails
- **Check**: Is the connection to the relay server still alive?
- **Check**: Is there network latency or firewall issues?

### 6. Unknown Message Type
- **Check**: Are you using a compatible client and server version?
- **Check**: Is the protocol version in hello message correct?

## Debugging Tips
- Run with `--verbose` to enable detailed logging.
- Check logs for error codes and messages.
- Use `openssl s_client` to debug TLS connections.
- Validate JWT tokens with [jwt.io](https://jwt.io/).
- Check relay server logs for more details.

## Getting Help
- Review the README and docs/README.md for configuration and usage.
- Open an issue on GitHub with logs and config details (do not include secrets).
- For security issues, contact the security contact in the main README. 