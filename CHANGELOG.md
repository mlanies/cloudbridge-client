# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Enhanced documentation with comprehensive guides
- Improved error handling and logging
- Better test coverage and architecture

### Changed
- Updated configuration structure for better flexibility
- Improved TLS 1.3 implementation
- Enhanced rate limiting algorithms

### Fixed
- Various bug fixes and performance improvements

## [1.1.1] - 2025-01-27

### Added
- Comprehensive technical specification compliance
- Improved test architecture with mock clients
- Enhanced documentation structure
- Better error handling and retry logic

### Changed
- Fixed tunnel tests to match current architecture
- Updated mock client implementation for proper interface compliance
- Improved test coverage and reliability
- Enhanced documentation with Russian and English versions

### Fixed
- Resolved merge conflicts in README.md
- Fixed compilation errors in test files
- Corrected interface implementations in tests
- Removed outdated and non-working test files

### Removed
- Deleted outdated metrics_test.go (not implemented)
- Removed unused imports and dead code
- Cleaned up deprecated test methods

## [1.1.0] - 2025-01-20

### Added
- Initial release with TLS 1.3 support
- JWT authentication implementation
- Cross-platform support (Windows, Linux, macOS)
- Comprehensive error handling
- Rate limiting and retry logic
- Heartbeat mechanism
- Tunnel management system

### Security
- Enforced TLS 1.3 with secure cipher suites
- JWT token validation with HMAC and RSA support
- Certificate validation and SNI support
- Rate limiting per user (JWT subject)

### Features
- YAML configuration with environment variable support
- Service installation and management
- Keycloak integration (OpenID Connect)
- Comprehensive logging and monitoring
- Cross-platform binary distribution

## [1.0.0] - 2025-01-15

### Added
- Initial project structure
- Basic client implementation
- Core protocol support
- Documentation framework

---

## Version History

- **v1.1.1**: Technical specification compliance, test fixes, documentation improvements
- **v1.1.0**: Full feature implementation with security and cross-platform support
- **v1.0.0**: Initial project setup and basic functionality

## Migration Guide

### From v1.0.x to v1.1.x
- Update configuration format to match new YAML structure
- Review authentication settings for JWT/Keycloak compatibility
- Test TLS 1.3 connectivity with your relay server
- Update service installation commands if using systemd integration

### From v1.1.0 to v1.1.1
- No breaking changes
- Improved test reliability and coverage
- Enhanced documentation and troubleshooting guides
- Better error handling and logging

## Support

For support and questions:
- Create an issue on GitHub
- Check the documentation in `/docs` directory
- Review configuration examples
- Consult troubleshooting guide 