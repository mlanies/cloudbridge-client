# CloudBridge Relay Client Documentation v2.0

## Overview

CloudBridge Relay Client v2.0 is a secure, cross-platform client for connecting to CloudBridge Relay servers with multi-tenancy support, Prometheus metrics, and performance optimization. It supports TLS 1.3, JWT/Keycloak authentication, tunnel management, heartbeat, rate limiting, and robust error handling.

---

## Architecture Overview

- **ConnectionManager**: Handles secure TLS 1.3 connections to the relay server.
- **AuthenticationManager**: Supports JWT (HS256) and optional Keycloak (OpenID Connect) authentication with tenant_id extraction.
- **TunnelManager**: Manages tunnel creation, validation, lifecycle, buffer management, and per-tunnel statistics.
- **HeartbeatManager**: Maintains connection health with periodic heartbeat messages.
- **ErrorHandler**: Centralized error handling and retry logic with exponential backoff, including multi-tenancy error codes.
- **Config**: YAML-based configuration, overridable via environment variables and CLI.
- **Metrics**: Prometheus-compatible metrics server for monitoring performance and tenant activity.
- **PerformanceOptimizer**: Automatic runtime optimization for high throughput or low latency workloads.

---

## New Features in v2.0

### Multi-Tenancy Support
- JWT tokens now require `tenant_id` claim
- All tunnel operations include tenant_id
- Per-tenant resource limits and monitoring
- Tenant isolation and IP filtering support

### Enhanced TCP Proxy
- Buffer management with connection pooling
- Per-tunnel statistics and monitoring
- Optimized data transfer with configurable buffer sizes

### Prometheus Metrics
- Real-time performance monitoring
- Per-tenant and per-tunnel statistics
- Buffer usage and connection metrics
- Error rate tracking

### Performance Optimization
- Automatic runtime optimization
- High throughput and low latency modes
- Memory ballast for better GC performance
- Configurable garbage collection settings

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

### Service Installation
```bash
# Linux/macOS
sudo cloudbridge-client service install <jwt-token>

# Windows
cloudbridge-client.exe service install <jwt-token>
```

---

## Configuration Reference

See `config.yaml` for a full example. All options can be set via environment variables (prefix `CLOUDBRIDGE_`).

### Core Settings
- **relay.host**: Relay server hostname
- **relay.port**: Relay server port
- **relay.tls.enabled**: Enforce TLS (must be true)
- **relay.tls.min_version**: Only "1.3" supported
- **relay.tls.verify_cert**: Enable certificate validation
- **relay.tls.ca_cert**: Path to CA certificate
- **auth.type**: "jwt" or "keycloak"
- **auth.secret**: JWT secret (for HS256)
- **auth.keycloak.enabled**: Enable Keycloak integration

### New v2.0 Settings
- **metrics.enabled**: Enable Prometheus metrics
- **metrics.prometheus_port**: Port for metrics endpoint (default: 9090)
- **metrics.tenant_metrics**: Enable per-tenant metrics
- **metrics.buffer_metrics**: Enable buffer pool metrics
- **performance.enabled**: Enable performance optimization
- **performance.optimization_mode**: "high_throughput" or "low_latency"
- **performance.gc_percent**: Garbage collection percentage
- **performance.memory_ballast**: Enable memory ballast

### Rate Limiting
- **rate_limiting.enabled**: Enable rate limiting
- **rate_limiting.max_retries**: Max retry attempts
- **rate_limiting.backoff_multiplier**: Exponential backoff multiplier
- **rate_limiting.max_backoff**: Max backoff duration

---

## Security Considerations

- **TLS 1.3 enforced**: Only secure cipher suites allowed
- **Certificate validation**: Strict, with optional CA pinning
- **JWT**: Only HS256 supported, with claim validation (`sub` and `tenant_id` required)
- **Keycloak**: OpenID Connect, automatic JWKS update, role/permission checks
- **Rate limiting**: Per-user (JWT subject) and per-tenant (tenant_id), exponential backoff, logging
- **Token storage**: Never log or persist tokens insecurely
- **Audit**: All operations are logged for audit purposes
- **Metrics security**: Prometheus endpoint should be restricted to internal network

---

## Error Handling & Troubleshooting

### Core Errors
- **invalid_token**: Check JWT validity, signature, expiration, and tenant_id
- **rate_limit_exceeded**: Too many requests; client will retry with backoff
- **connection_limit_reached**: Too many concurrent connections
- **server_unavailable**: Server is down or unreachable
- **invalid_tunnel_info**: Check tunnel parameters and tenant_id

### New v2.0 Errors
- **tenant_not_found**: tenant_id not found or not authorized
- **tenant_limit_exceeded**: Tenant has reached resource limits
- **ip_not_allowed**: Client IP not in allowed range for tenant
- **buffer_pool_exhausted**: Buffer pool exhausted, try later
- **data_transfer_failed**: Data transfer error

### Troubleshooting Steps
- Enable verbose logging (`--verbose`)
- Check relay server logs
- Validate TLS certificates and CA
- Ensure JWT secret matches relay server and includes tenant_id
- For Keycloak, check realm, client_id, and JWKS endpoint
- Monitor Prometheus metrics for performance insights

---

## Testing

### Unit Tests
```bash
go test ./...
```

### Integration Tests
- Requires relay-server binary (not included in this repository)
- Tests full connection cycle, TLS handshake, authentication, tunnel creation
- Can be mocked for CI/CD purposes

### Test Coverage
- Authentication (JWT, Keycloak, tenant_id extraction)
- Tunnel creation and validation (with tenant_id)
- Error handling (all error codes including multi-tenancy)
- Rate limiting and retry logic (per-tenant)
- Heartbeat manager
- BufferManager and performance optimization
- Prometheus metrics registration

---

## Deployment & Support

- See README for build and deployment instructions
- For issues, use the GitHub issue tracker
- For security concerns, contact the security contact listed in the README
- Monitor Prometheus metrics for operational insights

---

## Recent Updates (v2.0)

- Added multi-tenancy support with tenant_id in JWT tokens
- Implemented Prometheus metrics for real-time monitoring
- Enhanced TCP proxy with buffer management and statistics
- Added performance optimization with runtime tuning
- Extended error handling for multi-tenancy scenarios
- Updated all tests to support new features
- Improved documentation and configuration examples

---

## API & Protocol

All messages are JSON, UTF-8 encoded, no compression.

### Example: Hello
```json
{"type": "hello", "version": "1.0", "features": ["tls", "heartbeat", "tunnel_info", "multi_tenancy", "metrics"]}
```

### Example: Auth
```json
{"type": "auth", "token": "<jwt-with-tenant_id>"}
```

### Example: Tunnel
```json
{"type": "tunnel_info", "tunnel_id": "tunnel_001", "tenant_id": "tenant-001", "local_port": 3389, "remote_host": "192.168.1.100", "remote_port": 3389}
```

### Example: Heartbeat
```json
{"type": "heartbeat"}
```

### Example: Error
```json
{"type": "error", "code": "tenant_limit_exceeded", "message": "Tunnel limit exceeded for tenant"}
``` 