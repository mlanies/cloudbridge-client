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
go test ./pkg/relay
```

### Test Results (v1.1.1)
All unit tests pass successfully:
- ✅ Authentication tests (JWT validation, expired tokens, invalid signatures)
- ✅ Tunnel management tests (creation, validation, concurrency)
- ✅ Error handling tests
- ✅ Rate limiting tests

---

## Integration Testing

### Current Status
The integration test requires a relay-server binary which is not included in this repository. The test is located in `test/integration_test.go`.

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

### Mocking for CI/CD
For continuous integration, consider:
- Mocking the relay server responses
- Using a test relay server container
- Implementing stub responses for protocol testing

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

## Test Architecture

### Mock Client Implementation
Tests use a mock client that implements the `interfaces.ClientInterface`:

```go
type mockClient struct{}

func (m *mockClient) IsConnected() bool { return true }
func (m *mockClient) SendHeartbeat() error { return nil }
func (m *mockClient) GetConfig() *types.Config { return nil }
func (m *mockClient) GetClientID() string { return "mock-client" }
```

### Tunnel Manager Testing
Tunnel tests directly use `tunnel.Manager` with mock client:
- Tests tunnel registration and unregistration
- Validates tunnel parameters
- Tests concurrent tunnel operations
- Verifies tunnel lifecycle management

---

## Troubleshooting Failed Tests
- Use `--verbose` flag for detailed logs
- Check logs for error codes and stack traces
- Review `docs/TROUBLESHOOTING.md` for common issues
- Ensure all dependencies are properly installed

---

## Continuous Integration
- Integrate `go test ./...` into your CI pipeline (GitHub Actions, GitLab CI, etc.)
- Fail builds on test failures
- Consider adding test coverage reporting
- Implement integration test mocking for reliable CI/CD

---

## Test Coverage Goals
- Authentication: 100% (JWT validation, Keycloak integration)
- Tunnel Management: 100% (creation, validation, lifecycle)
- Error Handling: 100% (all error codes and retry logic)
- Rate Limiting: 100% (backoff strategies, limits)
- Heartbeat: 100% (connection health monitoring)

---

## Reporting Issues
- Open issues on GitHub with test logs and environment details (do not include secrets)
- Include Go version and platform information
- Provide minimal reproduction steps
- Attach relevant configuration files (without sensitive data) 