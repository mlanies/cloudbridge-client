# CloudBridge Relay Client API & Protocol

All messages are JSON, UTF-8 encoded, no compression.

## Message Types

### 1. Hello
- **Client → Server**
```json
{
  "type": "hello",
  "version": "1.0",
  "features": ["tls", "heartbeat", "tunnel_info"]
}
```
- **Server → Client**
```json
{
  "type": "hello_response",
  "version": "1.0",
  "features": ["tls", "heartbeat", "tunnel_info"]
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
  "client_id": "user123"
}
```

### 3. Tunnel Management
- **Client → Server**
```json
{
  "type": "tunnel_info",
  "tunnel_id": "tunnel_001",
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
  "code": "rate_limit_exceeded",
  "message": "Rate limit exceeded for user"
}
```

## Error Codes
- `invalid_token` — Invalid or expired JWT token
- `rate_limit_exceeded` — Rate limit exceeded
- `connection_limit_reached` — Connection limit reached
- `server_unavailable` — Server unavailable
- `invalid_tunnel_info` — Invalid tunnel info
- `unknown_message_type` — Unknown message type

## Notes
- All fields are required unless otherwise specified.
- All messages must be valid UTF-8 JSON.
- No message compression is used.
- All connections must use TLS 1.3. 