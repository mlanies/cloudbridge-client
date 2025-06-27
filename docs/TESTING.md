# Testing Guide: CloudBridge Relay Client v2.0

## Unit Testing

### Running All Unit Tests
```bash
go test ./...
```

### What to Test
- Authentication (JWT, Keycloak, tenant_id extraction)
- Tunnel creation and validation (с tenant_id)
- Error handling (all error codes, включая multi-tenancy)
- Rate limiting and retry logic
- Heartbeat manager
- BufferManager and performance optimizer
- Prometheus metrics registration (unit)

### Example: Run Tests for a Specific Package
```bash
go test ./pkg/auth
go test ./pkg/relay
```

### Test Results (v2.0)
All unit tests pass successfully:
- ✅ Authentication tests (JWT validation, expired tokens, invalid signatures, tenant_id extraction)
- ✅ Tunnel management tests (creation, validation, concurrency, tenant_id)
- ✅ Error handling tests (multi-tenancy, performance)
- ✅ Rate limiting tests
- ✅ BufferManager and performance tests

---

## Integration Testing

### Current Status
The integration test requires a relay-server binary which is not included in this repository. The test is located in `test/integration_test.go`.

### Full Connection Cycle
- Start a test relay server (with TLS 1.3)
- Run the client with a valid JWT token (with tenant_id)
- Verify connection, authentication, tunnel creation, heartbeat, metrics

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
- Test with missing or invalid `sub` or `tenant_id` claim

### Rate Limiting
- Simulate burst requests to trigger rate limiting
- Verify exponential backoff and retry logic
- Test per-tenant limits

### Penetration Testing
- Use tools like `nmap`, `sslscan`, and custom scripts to test for vulnerabilities
- Attempt replay, DoS, and protocol fuzzing attacks
- Test tenant isolation and metrics endpoint security

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
func (m *mockClient) GetTenantID() string { return "mock-tenant" }
```

### Tunnel Manager Testing
Tunnel tests directly use `tunnel.Manager` with mock client:
- Tests tunnel registration and unregistration
- Validates tunnel parameters (including tenant_id)
- Tests concurrent tunnel operations
- Verifies tunnel lifecycle management

### Metrics Testing
- Unit-test Prometheus metrics registration (no panic)
- Optionally, run client and check `/metrics` endpoint for expected stats

### Performance Testing
- Test BufferManager under load (concurrent tunnel creation)
- Test optimizer settings (GOMAXPROCS, GC)

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
- Authentication: 100% (JWT validation, Keycloak integration, tenant_id extraction)
- Tunnel Management: 100% (creation, validation, lifecycle, tenant_id)
- Error Handling: 100% (all error codes and retry logic)
- Rate Limiting: 100% (backoff strategies, limits, per-tenant)
- Heartbeat: 100% (connection health monitoring)
- BufferManager: 100% (buffer pool, exhaustion)
- Metrics: 100% (registration, endpoint)

---

## Reporting Issues
- Open issues on GitHub with test logs and environment details (do not include secrets)
- Include Go version and platform information
- Provide minimal reproduction steps
- Attach relevant configuration files (without sensitive data) 