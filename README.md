# CloudBridge Relay Client

A cross-platform Go implementation of a client for the CloudBridge Relay service. This client allows you to establish secure tunnels through the CloudBridge Relay server using JWT authentication.

## Features

- Cross-platform support (Windows, Linux, macOS)
- Secure TLS connection support
- JWT token authentication
- Automatic reconnection handling
- Graceful shutdown
- Configurable local port tunneling
- Comprehensive error handling

## Installation

### Using Go Install

```bash
go install github.com/2gc-dev/cloudbridge-client/cmd/cloudbridge-client@latest
```

### Pre-built Binaries

Download the appropriate binary for your platform from the [releases page](https://github.com/2gc-dev/cloudbridge-client/releases).

## Usage

### Basic Usage

```bash
cloudbridge-client --token "your-jwt-token"
```

This will connect to the default relay server (edge.2gc.ru:8080) with TLS enabled.

### Custom Configuration

```bash
cloudbridge-client \
  --host custom-relay-server.com \
  --port 8080 \
  --token "your-jwt-token" \
  --local-port 3389
```

### With Custom TLS Certificates

```bash
cloudbridge-client \
  --cert client.crt \
  --key client.key \
  --ca ca.crt \
  --token "your-jwt-token"
```

### Command Line Options

- `--tls`: Enable TLS connection (default: true)
- `--cert`: Path to TLS certificate file
- `--key`: Path to TLS key file
- `--ca`: Path to CA certificate file
- `--host`: Relay server hostname (default: edge.2gc.ru)
- `--port`: Relay server port (default: 8080)
- `--token`: JWT token for authentication (required)
- `--local-port`: Local port to tunnel (default: 3389)
- `--reconnect-delay`: Reconnection delay in seconds (default: 5)
- `--max-retries`: Maximum number of reconnection attempts (default: 3)

## Platform Support

The client is tested and supported on:

- Windows (x86_64, ARM64)
- Linux (x86_64, ARM64)
- macOS (x86_64, ARM64)

## Security Considerations

1. TLS is enabled by default for all connections
2. Keep your JWT tokens secure and rotate them regularly
3. Use strong TLS certificates and keys
4. Validate server certificates using a trusted CA

## Error Handling

The client handles various error conditions:

- Invalid JWT tokens
- Rate limiting
- Connection limits
- Server unavailability
- TLS certificate issues

## Building from Source

```bash
git clone https://github.com/2gc-dev/cloudbridge-client.git
cd cloudbridge-client
go build -o cloudbridge-client ./cmd/cloudbridge-client
```

### Cross-Platform Build

To build for multiple platforms:

```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o cloudbridge-client.exe ./cmd/cloudbridge-client

# Linux
GOOS=linux GOARCH=amd64 go build -o cloudbridge-client ./cmd/cloudbridge-client

# macOS
GOOS=darwin GOARCH=amd64 go build -o cloudbridge-client ./cmd/cloudbridge-client
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 