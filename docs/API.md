# CloudBridge Relay Client API & Protocol (v2.0)

All messages are JSON, UTF-8 encoded, no compression.

## Message Types

### 1. Hello
- **Client → Server**
```json
{
  "type": "hello",
  "version": "1.0",
  "features": ["tls", "heartbeat", "tunnel_info", "multi_tenancy", "metrics"]
}
```
- **Server → Client**
```json
{
  "type": "hello_response",
  "version": "1.0",
  "features": ["tls", "heartbeat", "tunnel_info", "multi_tenancy", "metrics"]
}
```

### 2. Authentication
- **Client → Server**
```json
{
  "type": "auth",
  "token": "<jwt>"
}
```
- **Server → Client**
```json
{
  "type": "auth_response",
  "status": "ok",
  "client_id": "user123",
  "tenant_id": "tenant-001"
}
```

### 3. Tunnel Management
- **Client → Server**
```json
{
  "type": "tunnel_info",
  "tunnel_id": "tunnel_001",
  "tenant_id": "tenant-001",
  "local_port": 3389,
  "remote_host": "192.168.1.100",
  "remote_port": 3389
}
```
- **Server → Client**
```json
{
  "type": "tunnel_response",
  "status": "ok",
  "tunnel_id": "tunnel_001"
}
```

### 4. Heartbeat
- **Client → Server**
```json
{
  "type": "heartbeat"
}
```
- **Server → Client**
```json
{
  "type": "heartbeat_response"
}
```

### 5. Error
- **Server → Client**
```json
{
  "type": "error",
  "code": "tenant_limit_exceeded",
  "message": "Tunnel limit exceeded for tenant"
}
```

## Error Codes
- `invalid_token` — Invalid or expired JWT token
- `rate_limit_exceeded` — Rate limit exceeded
- `connection_limit_reached` — Connection limit reached
- `server_unavailable` — Server unavailable
- `invalid_tunnel_info` — Invalid tunnel info
- `unknown_message_type` — Unknown message type
- `tenant_limit_exceeded` — Tunnel or connection limit for tenant exceeded
- `tenant_not_found` — Tenant not found or not authorized
- `ip_not_allowed` — IP address not allowed for this tenant
- `buffer_pool_exhausted` — Buffer pool exhausted, try later
- `data_transfer_failed` — Data transfer error

## Notes
- All fields are required unless otherwise specified.
- All messages must be valid UTF-8 JSON.
- No message compression is used.
- All connections must use TLS 1.3.
- For multi-tenancy, tenant_id is required in JWT and all tunnel-related messages.
- Metrics are available via Prometheus endpoint if enabled in config. 