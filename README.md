# CloudBridge Relay Client

A cross-platform Go implementation of a client for the CloudBridge Relay service. This client implements the complete protocol specification with TLS 1.3 support, JWT authentication, and comprehensive error handling.

## Features

- **TLS 1.3 Support**: Enforced TLS 1.3 with secure cipher suites
- **JWT Authentication**: Full JWT token validation with HMAC and RSA support
- **Keycloak Integration**: Optional OpenID Connect integration
- **Cross-platform**: Windows, Linux, macOS (x86_64, ARM64)
- **Rate Limiting**: Built-in rate limiting with exponential backoff
- **Heartbeat**: Automatic connection health monitoring
- **Tunnel Management**: Full tunnel lifecycle management
- **Error Handling**: Comprehensive error handling and retry logic
- **Configuration**: Flexible YAML configuration with environment variable support

## Protocol Support

This client implements the complete CloudBridge Relay protocol:

- **Hello/Hello Response**: Protocol version negotiation
- **Auth/Auth Response**: JWT-based authentication
- **Tunnel Info/Tunnel Response**: Tunnel creation and management
- **Heartbeat/Heartbeat Response**: Connection health monitoring
- **Error Messages**: Standardized error handling

## Installation

### Using Go Install

```bash
go install github.com/2gc-dev/cloudbridge-client/cmd/cloudbridge-client@latest
```

### Pre-built Binaries

Download the appropriate binary for your platform from the [releases page](https://github.com/2gc-dev/cloudbridge-client/releases).

### Building from Source

```bash
git clone https://github.com/2gc-dev/cloudbridge-client.git
cd cloudbridge-client
go build -o cloudbridge-client ./cmd/cloudbridge-client
```

## Quick Start

### Basic Usage

```bash
cloudbridge-client --token "your-jwt-token"
```

This will connect to the default relay server (edge.2gc.ru:8080) with TLS enabled.

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

## Configuration

The client supports configuration via YAML files and environment variables.

### Configuration File (config.yaml)

```yaml
relay:
  host: "edge.2gc.ru"
  port: 8080
  timeout: "30s"
  tls:
    enabled: true
    min_version: "1.3"
    verify_cert: true
    ca_cert: "/path/to/ca.pem"
    client_cert: "/path/to/client.crt"
    client_key: "/path/to/client.key"

auth:
  type: "jwt"
  secret: "your-jwt-secret"
  keycloak:
    enabled: false
    server_url: "https://keycloak.example.com"
    realm: "cloudbridge"
    client_id: "relay-client"

rate_limiting:
  enabled: true
  max_retries: 3
  backoff_multiplier: 2.0
  max_backoff: "30s"

logging:
  level: "info"
  format: "json"
  output: "stdout"
```

### Environment Variables

All configuration options can be set via environment variables with the `CLOUDBRIDGE_` prefix:

```bash
export CLOUDBRIDGE_RELAY_HOST="edge.2gc.ru"
export CLOUDBRIDGE_RELAY_PORT="8080"
export CLOUDBRIDGE_AUTH_SECRET="your-jwt-secret"
```

### Command Line Options

- `--config, -c`: Configuration file path
- `--token, -t`: JWT token for authentication (required)
- `--tunnel-id, -i`: Tunnel ID (default: tunnel_001)
- `--local-port, -l`: Local port to bind (default: 3389)
- `--remote-host, -r`: Remote host (default: 192.168.1.100)
- `--remote-port, -p`: Remote port (default: 3389)
- `--verbose, -v`: Enable verbose logging

## Security Features

### TLS 1.3

- Enforced TLS 1.3 minimum version
- Secure cipher suites only:
  - `TLS_AES_256_GCM_SHA384`
  - `TLS_CHACHA20_POLY1305_SHA256`
  - `TLS_AES_128_GCM_SHA256`
- Certificate validation
- SNI support

### JWT Authentication

- HMAC-SHA256 support
- RSA signature validation
- Token expiration checking
- Subject extraction for rate limiting

### Keycloak Integration

- OpenID Connect support
- Automatic JWKS fetching
- Token validation
- Role-based access control

## Error Handling

The client handles all standard relay errors:

- `invalid_token`: Invalid or expired JWT token
- `rate_limit_exceeded`: Rate limiting with exponential backoff
- `connection_limit_reached`: Connection limit exceeded
- `server_unavailable`: Server unavailability with retry
- `invalid_tunnel_info`: Invalid tunnel configuration
- `unknown_message_type`: Protocol errors

## Rate Limiting

Built-in rate limiting with configurable parameters:

- Per-user rate limiting (based on JWT subject)
- Exponential backoff retry strategy
- Configurable maximum retries
- Maximum backoff limits

## Heartbeat

Automatic connection health monitoring:

- Configurable heartbeat interval (default: 30s)
- Failure detection and handling
- Automatic reconnection on failures
- Heartbeat statistics

## Platform Support

Tested and supported on:

- **Windows**: x86_64, ARM64
- **Linux**: x86_64, ARM64
- **macOS**: x86_64, ARM64

## Development

### Building for Multiple Platforms

```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o cloudbridge-client.exe ./cmd/cloudbridge-client

# Linux
GOOS=linux GOARCH=amd64 go build -o cloudbridge-client ./cmd/cloudbridge-client

# macOS
GOOS=darwin GOARCH=amd64 go build -o cloudbridge-client ./cmd/cloudbridge-client
```

### Running Tests

```bash
go test ./...
```

### Code Structure

```
pkg/
├── auth/          # Authentication management
├── config/        # Configuration handling
├── errors/        # Error handling and retry logic
├── heartbeat/     # Heartbeat management
├── relay/         # Main relay client
└── tunnel/        # Tunnel management

cmd/
└── cloudbridge-client/  # Main application
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support and questions:

- Create an issue on GitHub
- Check the documentation
- Review the configuration examples

## Changelog

### v1.0.0
- Initial release
- TLS 1.3 support
- JWT authentication
- Cross-platform support
- Comprehensive error handling
- Rate limiting and retry logic
- Heartbeat mechanism
- Tunnel management 