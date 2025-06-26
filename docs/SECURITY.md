# Security Considerations for CloudBridge Relay Client

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

## Authentication
- **JWT**: Only HS256 supported, secret must match relay server
- **Keycloak**: OpenID Connect, automatic JWKS update, role/permission checks
- **Claims**: `sub` (subject) required for user identification
- **Token validation**: signature, expiration, issuer
- **Never log or persist tokens in plaintext**

## Rate Limiting
- Per-user (by JWT subject)
- Default: 100 requests/sec, burst 200
- Exponential backoff, max 3 retries
- All rate limit violations are logged

## Logging & Audit
- All operations are logged (level configurable)
- Sensitive data (tokens, secrets) never logged
- Audit logs should be protected and regularly reviewed

## Error Handling
- All protocol errors are handled and logged
- Exponential backoff for retryable errors
- Graceful shutdown on fatal errors

## Secure Deployment
- Store config files and secrets securely (use environment variables for secrets if possible)
- Restrict access to config.yaml and logs
- Regularly update dependencies and perform security audits

## Penetration Testing
- Regularly test for vulnerabilities (TLS, JWT, replay, DoS)
- Validate server certificates and JWTs in all test cases

## Contact
- For security issues, contact the security contact listed in the main README. 