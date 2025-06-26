# Testing Guide: CloudBridge Relay Client

## Unit Testing

### Running All Unit Tests
```bash
go test ./...
```

### What to Test
- Authentication (JWT, Keycloak)
- Tunnel creation and validation
- Error handling (all error codes)
- Rate limiting and retry logic
- Heartbeat manager

### Example: Run Tests for a Specific Package
```bash
go test ./pkg/auth
```

---

## Integration Testing

### Full Connection Cycle
- Start a test relay server (with TLS 1.3)
- Run the client with a valid JWT token
- Verify connection, authentication, tunnel creation, heartbeat

### TLS 1.3 Handshake
- Use `openssl s_client -connect relay.example.com:8080 -tls1_3` to verify server
- Run client and check for successful handshake

### Real Relay Server
- Deploy relay server in test environment
- Run client end-to-end with real tokens and tunnels

---

## Security Testing

### Certificate Validation
- Test with valid and invalid CA certificates
- Test with expired or revoked server certificates

### JWT Validation
- Test with valid, expired, and tampered tokens
- Test with missing or invalid `sub` claim

### Rate Limiting
- Simulate burst requests to trigger rate limiting
- Verify exponential backoff and retry logic

### Penetration Testing
- Use tools like `nmap`, `sslscan`, and custom scripts to test for vulnerabilities
- Attempt replay, DoS, and protocol fuzzing attacks

---

## Troubleshooting Failed Tests
- Use `--verbose` flag for detailed logs
- Check logs for error codes and stack traces
- Review `docs/TROUBLESHOOTING.md` for common issues

---

## Continuous Integration
- Integrate `go test ./...` into your CI pipeline (GitHub Actions, GitLab CI, etc.)
- Fail builds on test failures

---

## Reporting Issues
- Open issues on GitHub with test logs and environment details (do not include secrets) 