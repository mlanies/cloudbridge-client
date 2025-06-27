# Security Considerations for CloudBridge Relay Client v2.0

## TLS
- Enforced TLS 1.3 only
- Only secure cipher suites allowed:
  - TLS_AES_256_GCM_SHA384
  - TLS_CHACHA20_POLY1305_SHA256
  - TLS_AES_128_GCM_SHA256
- Certificate validation is mandatory
- SNI (Server Name Indication) supported
- Optional CA pinning via config
- HSTS and OCSP stapling recommended on server side

## Authentication & Multi-Tenancy
- **JWT**: Only HS256 supported, secret must match relay server
- **Keycloak**: OpenID Connect, automatic JWKS update, role/permission checks
- **Claims**: `sub` (subject) required for user identification, `tenant_id` required for multi-tenancy
- **Token validation**: signature, expiration, issuer, tenant_id
- **Never log or persist tokens in plaintext**
- **tenant_id** must be validated on server and not be user-controlled

## Rate Limiting
- Per-user (by JWT subject) and per-tenant (by tenant_id)
- Default: 100 requests/sec, burst 200
- Exponential backoff, max 3 retries
- All rate limit violations are logged
- New error codes: `tenant_limit_exceeded`, `tenant_not_found`, `ip_not_allowed`

## Logging & Audit
- All operations are logged (level configurable)
- Sensitive data (tokens, secrets) never logged
- Audit logs should be protected and regularly reviewed
- tenant_id всегда логируется для всех операций в multi-tenancy режиме

## Error Handling
- All protocol errors are handled and logged
- Exponential backoff for retryable errors
- Graceful shutdown on fatal errors
- Новые error code для multi-tenancy и performance (см. API.md)

## Secure Deployment
- Store config files and secrets securely (use environment variables for secrets if possible)
- Restrict access to config.yaml and logs
- Regularly update dependencies and perform security audits
- Restrict Prometheus /metrics endpoint to internal network only (use firewall or listen on localhost)

## Prometheus & Monitoring
- Prometheus endpoint (`/metrics`) exposes operational metrics, including per-tenant statistics
- Do not expose /metrics to the public internet
- Monitor for unusual activity per tenant (spikes, errors, resource exhaustion)

## Penetration Testing
- Regularly test for vulnerabilities (TLS, JWT, replay, DoS, tenant isolation)
- Validate server certificates and JWTs in all test cases
- Test for tenant boundary violations and metrics leaks

## Contact
- For security issues, contact the security contact listed in the main README. 