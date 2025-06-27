# Architecture Overview: CloudBridge Relay Client v2.0

## Main Components

- **ConnectionManager**: Establishes and maintains secure TLS 1.3 connections to the relay server.
- **AuthenticationManager**: Handles JWT and Keycloak authentication, token validation, and claim extraction (включая tenant_id для multi-tenancy).
- **TunnelManager**: Manages tunnel creation, validation, lifecycle (local/remote port mapping, proxying), buffer management, and per-tunnel statistics.
- **HeartbeatManager**: Periodically sends heartbeat messages to monitor connection health and trigger reconnection if needed.
- **ErrorHandler**: Centralized error handling, retry logic, and exponential backoff for transient errors, включая новые error code для multi-tenancy и performance.
- **Config**: Loads and validates configuration from YAML, environment variables, and CLI flags (включая секции metrics и performance).
- **Logging**: Structured logging with configurable level and format.
- **Metrics**: Prometheus-compatible metrics server, собирает статистику по трафику, соединениям, буферам, tenant-метрикам.
- **PerformanceOptimizer**: Автоматическая оптимизация runtime (GOMAXPROCS, GC, memory ballast) для throughput/latency.

## Data Flow

1. **Startup**: Load config → parse CLI/env → validate
2. **Connect**: Establish TLS 1.3 connection to relay
3. **Hello**: Exchange hello/hello_response messages (protocol negotiation)
4. **Authenticate**: Send JWT/Keycloak token (с tenant_id), receive auth_response
5. **Tunnel**: Send tunnel_info (с tenant_id), receive tunnel_response, start proxy (buffered)
6. **Heartbeat**: Periodically send heartbeat, handle heartbeat_response
7. **Metrics**: Все события и трафик отражаются в Prometheus-метриках
8. **Performance**: Оптимизация runtime на старте и в процессе работы
9. **Error Handling**: On error, apply retry/backoff or shutdown gracefully

## Component Interaction Diagram

```mermaid
graph TD
  A[Config Loader] --> B[ConnectionManager]
  B --> C[AuthenticationManager]
  C --> D[TunnelManager]
  B --> E[HeartbeatManager]
  B --> F[ErrorHandler]
  F --> B
  F --> C
  F --> D
  F --> E
  B --> G[Logging]
  C --> G
  D --> G
  E --> G
  F --> G
  D --> H[Metrics]
  B --> I[PerformanceOptimizer]
```

## Extensibility
- New authentication methods can be added via AuthenticationManager
- Additional tunnel types or protocols can be added to TunnelManager
- Logging and monitoring can be integrated via Logging and Metrics components
- Performance tuning can be extended via PerformanceOptimizer

## Security Boundaries
- All network traffic is encrypted (TLS 1.3)
- Tokens and secrets are never logged
- All errors and retries are logged for audit
- tenant_id обязателен для multi-tenancy и валидируется на сервере
- Prometheus endpoint рекомендуется ограничивать внутренней сетью

## See Also
- `docs/README.md` for usage
- `docs/API.md` for protocol details
- `docs/SECURITY.md` for security model 