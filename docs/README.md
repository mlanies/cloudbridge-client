# CloudBridge Relay Client Documentation

## Overview

CloudBridge Relay Client is a secure, cross-platform client for connecting to CloudBridge Relay servers. It supports TLS 1.3, JWT/Keycloak authentication, tunnel management, heartbeat, rate limiting, and robust error handling.

---

## Architecture Overview

- **ConnectionManager**: Handles secure TLS 1.3 connections to the relay server.
- **AuthenticationManager**: Supports JWT (HS256) and optional Keycloak (OpenID Connect) authentication.
- **TunnelManager**: Manages tunnel creation, validation, and lifecycle.
- **HeartbeatManager**: Maintains connection health with periodic heartbeat messages.
- **ErrorHandler**: Centralized error handling and retry logic with exponential backoff.
- **Config**: YAML-based configuration, overridable via environment variables and CLI.

---

## Usage Examples

### Basic Connection
```bash
cloudbridge-client --token "your-jwt-token"
```

### With Configuration File
```bash
cloudbridge-client --config config.yaml --token "your-jwt-token"
```

### Custom Tunnel
```bash
cloudbridge-client \
  --token "your-jwt-token" \
  --tunnel-id "my-tunnel" \
  --local-port 3389 \
  --remote-host "192.168.1.100" \
  --remote-port 3389
```

---

## Configuration Reference

See `config.yaml` for a full example. All options can be set via environment variables (prefix `CLOUDBRIDGE_`).

- **relay.host**: Relay server hostname
- **relay.port**: Relay server port
- **relay.tls.enabled**: Enforce TLS (must be true)
- **relay.tls.min_version**: Only "1.3" supported
- **relay.tls.verify_cert**: Enable certificate validation
- **relay.tls.ca_cert**: Path to CA certificate
- **auth.type**: "jwt" or "keycloak"
- **auth.secret**: JWT secret (for HS256)
- **auth.keycloak.enabled**: Enable Keycloak integration
- **rate_limiting.enabled**: Enable rate limiting
- **rate_limiting.max_retries**: Max retry attempts
- **rate_limiting.backoff_multiplier**: Exponential backoff multiplier
- **rate_limiting.max_backoff**: Max backoff duration

---

## Security Considerations

- **TLS 1.3 enforced**: Only secure cipher suites allowed
- **Certificate validation**: Strict, with optional CA pinning
- **JWT**: Only HS256 supported, with claim validation (`sub` required)
- **Keycloak**: OpenID Connect, automatic JWKS update, role/permission checks
- **Rate limiting**: Per-user (JWT subject), exponential backoff, logging
- **Token storage**: Never log or persist tokens insecurely
- **Audit**: All operations are logged for audit purposes

---

## Error Handling & Troubleshooting

- **invalid_token**: Check JWT validity, signature, and expiration
- **rate_limit_exceeded**: Too many requests; client will retry with backoff
- **connection_limit_reached**: Too many concurrent connections
- **server_unavailable**: Server is down or unreachable
- **invalid_tunnel_info**: Check tunnel parameters
- **unknown_message_type**: Protocol mismatch or bug

### Troubleshooting Steps
- Enable verbose logging (`--verbose`)
- Check relay server logs
- Validate TLS certificates and CA
- Ensure JWT secret matches relay server
- For Keycloak, check realm, client_id, and JWKS endpoint

---

## Acceptance Criteria Checklist

- [x] TLS 1.3 enforced, strict cipher suites
- [x] JWT authentication with claim validation
- [x] Tunnel management (create, validate, proxy)
- [x] Heartbeat and connection health
- [x] Rate limiting and retry logic
- [x] Comprehensive error handling
- [x] Logging and audit
- [x] Configurable via YAML, env, CLI

---

## API & Protocol

All messages are JSON, UTF-8 encoded, no compression.

### Example: Hello
```json
{"type": "hello", "version": "1.0", "features": ["tls", "heartbeat", "tunnel_info"]}
```

### Example: Auth
```json
{"type": "auth", "token": "<jwt>"}
```

### Example: Tunnel
```json
{"type": "tunnel_info", "tunnel_id": "tunnel_001", "local_port": 3389, "remote_host": "192.168.1.100", "remote_port": 3389}
```

### Example: Heartbeat
```json
{"type": "heartbeat"}
```

### Example: Error
```json
{"type": "error", "code": "rate_limit_exceeded", "message": "Rate limit exceeded for user"}
```

---

## Testing

- Unit tests: authentication, tunnel, error handling, rate limiting
- Integration: full connection cycle, TLS handshake, real relay server
- Security: certificate validation, JWT validation, penetration testing

---

## Deployment & Support

- See README for build and deployment instructions
- For issues, use the GitHub issue tracker
- For security concerns, contact the security contact listed in the README 