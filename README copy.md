# CloudBridge Relay v2.0

CloudBridge Relay v2.0 — высокопроизводительный сервер туннелирования с поддержкой TCP proxy и multi-tenancy, обеспечивающий защищённый доступ к внутренним ресурсам с централизованной аутентификацией через Keycloak. Это альтернатива Cloudflare Argo Tunnel с полным контролем над инфраструктурой и данными.

## 🚀 Новые возможности v2.0

### Ключевые улучшения
- **🔗 TCP Proxy**: Полнофункциональный TCP прокси для двунаправленной передачи данных
- **🏢 Multi-tenancy**: Поддержка мультитенантности с изоляцией ресурсов и IP-фильтрацией
- **⚡ Go 1.21**: Обновление до Go 1.21.13 с современными зависимостями
- **🔐 Enhanced Security**: Улучшенная интеграция с Keycloak и Django fallback
- **📊 Comprehensive Testing**: Полный набор unit и integration тестов
- **🛠️ Tunnel Manager**: Централизованное управление активными туннелями

### Технические улучшения
- **Производительность**: Улучшенная производительность TCP proxy с буферизацией
- **Мониторинг**: Расширенные метрики по тенантам и туннелям
- **Безопасность**: Улучшенная изоляция ресурсов между тенантами
- **Масштабируемость**: Поддержка тысяч одновременных соединений
- **Надежность**: Graceful shutdown и автоматическая очистка ресурсов

## Основная концепция

CloudBridge Relay v2.0 позволяет безопасно туннелировать доступ к внутренним сервисам (AD/LDAP, RDP, базы данных) через публичный Relay Server с поддержкой multi-tenancy:

* Безопасный доступ к внутренним сервисам без VPN
* Централизованная аутентификация через Keycloak с Django fallback
* Полная изоляция между компаниями (мультитенантность)
* TCP proxy для эффективной передачи данных
* Контроль и аудит операций на стороне Keycloak и Django
* Поведенческий анализ и ML-детекция угроз
* Интеграция с SIEM и внешними системами безопасности

## Архитектура компонентов

### Ключевые компоненты

#### 🔐 Keycloak (или другой IdP)
* Централизованное хранилище пользователей, групп и ролей
* Поддержка OIDC, OAuth2, LDAP, MFA
* JWKS (RS256) валидация токенов на стороне Relay
* Realm = изолированный клиент/арендатор
* Группы и атрибуты → маппинг на доступ к туннелям/сервисам
* **Fallback на Django аутентификацию при недоступности**

#### ⚙️ Django API
* Управление компаниями, серверами, тарифами, сервисами
* Выдача серверных токенов (Server Client)
* Web UI для администраторов и клиентов
* REST API для внешних систем (SIEM, DLP, партнеры)
* Аудит, биллинг, роли и разграничение доступа
* **Управление multi-tenancy и лимитами ресурсов**

#### 🚇 Relay Server v2.0 (Go 1.21)
* **TCP Proxy** для прокладки туннелей между Server Client ↔ Desktop Client
* Валидация JWT токенов через Keycloak с Django fallback
* Rate limiting, heartbeat, TLS/mTLS
* **Multi-tenancy** с изоляцией туннелей по tenant
* **Tunnel Manager** для централизованного управления туннелями
* Интеграция с Prometheus (расширенные метрики)
* **Автоматическая очистка неактивных ресурсов**

#### 💻 Server Client (Go)
* Развёртывается на стороне защищённой инфраструктуры
* Запрашивает токен у Django → подключается к Relay
* Объявляет список локальных сервисов (LDAP, RDP)
* **Поддержка TCP proxy для эффективной передачи данных**
* **Multi-tenant аутентификация**

#### 🖥️ Desktop Client (Tauri или CLI)
* Выполняет OIDC login через Keycloak
* Получает JWT токен
* Подключается к Relay Server
* Автоматически подключается к нужному туннелю (RDP, SSH)
* Скачивает .rdp-файл или запускает RDP-клиент
* **Поддержка multi-tenancy и TCP proxy**

#### 📊 SIEM-модуль (Python + Elasticsearch + Kafka)
* Приём событий из Relay/Django (входы, туннели, ошибки)
* Обогащение: user-agent, IP, geolocation, threat-level
* Анализ в реальном времени + историческое хранилище
* Отправка alert-ов → Telegram/Email/API
* UI-дэшборды (Grafana или встроенные)
* **Мониторинг multi-tenant активности**

#### 🧠 ML-модуль
* Feature extraction: паттерны времени, частоты, роли, IP
* Обучение моделей: классификация/аномалия/кластеризация
* Реакция: mark-as-suspicious, disable token, alert
* Отчёты по пользовательским профилям риска
* Режим обучения + режим обнаружения
* **Анализ поведения по тенантам**

#### 🔌 Open API Gateway (REST + Webhooks)
* Открытое API для:
  * событий (входы, туннели, логи)
  * управления туннелями
  * просмотра статуса подключения
  * интеграции с DLP/SOAR
  * **управления multi-tenancy**
* Webhooks для событий: tunnel_opened, login_failed, alert_triggered

#### 🚀 Production Infrastructure
* **High Availability (HA)** - кластеризация Relay Server с балансировкой нагрузки
* **Auto Scaling** - автоматическое масштабирование на основе метрик
* **Edge Networks** - распределенная сеть edge-узлов для глобального покрытия
* **Load Balancing** - L4/L7 балансировка с health checks
* **Circuit Breaker** - защита от каскадных отказов
* **Multi-tenant isolation** - полная изоляция ресурсов между тенантами

#### 🔒 Advanced Security
* **mTLS (Mutual TLS)** - двусторонняя аутентификация между Relay и клиентами
* **Misuse Protection** - защита от злоупотреблений и атак
* **Rate Limiting** - многоуровневое ограничение скорости
* **DDoS Protection** - защита от распределенных атак
* **IP Whitelisting** - белые списки IP-адресов по тенантам
* **Session Management** - управление сессиями и их жизненным циклом
* **Tenant isolation** - полная изоляция ресурсов между клиентами

#### 📊 Performance & Monitoring
* **Load Testing** - тестирование производительности под нагрузкой
* **Capacity Planning** - планирование ресурсов
* **Performance Metrics** - детальные метрики производительности
* **Resource Optimization** - оптимизация использования ресурсов
* **Multi-tenant metrics** - метрики по тенантам и туннелям

### Общий поток взаимодействия v2.0

#### Server Client:

1. Администратор создаёт сервер через Django UI.
2. Server Client запрашивает токен через Django API с указанием tenant_id.
3. Django делает запрос к Keycloak от имени Server Client.
4. Keycloak возвращает подписанный JWT токен с tenant_id.
5. Server Client подключается к Relay Server с этим токеном.
6. **Relay создает TCP proxy для эффективной передачи данных.**

#### Relay Server v2.0:

1. Получает токен от Server Client с tenant_id.
2. Валидирует его через Keycloak JWKS (RS256) или Django fallback.
3. Проверяет лимиты ресурсов для тенанта.
4. Извлекает информацию о туннелях.
5. **Создаёт TCP proxy туннели и начинает пересылку трафика.**
6. **Отслеживает активность и автоматически очищает неактивные ресурсы.**

#### Desktop Client:

1. Пользователь аутентифицируется через Keycloak (включая LDAP federation).
2. Получает JWT токен с tenant_id.
3. Подключается к Relay Server и получает доступ к нужным туннелям по правам.
4. **Использует TCP proxy для эффективного доступа к сервисам.**

## Потоки данных v2.0

### 🔐 Аутентификация с Multi-tenancy
```
User → Keycloak (OIDC login) → Access Token (tenant_id) → Relay
Server Client → Django → Keycloak → Token (tenant_id) → Relay
Relay → Keycloak JWKS validation → Django fallback (if needed)
```

### 🔄 TCP Proxy Туннелирование
```
Desktop Client ↔ Relay TCP Proxy ↔ Server Client ↔ Internal Service
                (Port 10000-20000)  (Bidirectional)
```

### 📡 Метрики и события
```
Relay / Django → Prometheus → Grafana (multi-tenant dashboards)
Relay / Django → Kafka → SIEM → AlertManager
SIEM → Telegram, Email, API
```

### 🧠 Поведенческий анализ
```
SIEM Logs → ML-модель → Risk Score (per tenant)
ML-модель → Response Engine (отключение токена, предупреждение)
```

## 🔗 TCP Proxy Архитектура

### Модель TCP Proxy

```
┌─────────────────────────────────────────────────────────────────┐
│                    TCP Proxy Architecture                      │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐         │
│  │   Desktop   │    │   TCP       │    │   Server    │         │
│  │   Client    │    │   Proxy     │    │   Client    │         │
│  │             │    │             │    │             │         │
│  │ localhost:  │◄──►│ Port        │◄──►│ Internal    │         │
│  │ 12345       │    │ 12345       │    │ Service     │         │
│  │             │    │ Buffer:     │    │ 192.168.1.  │         │
│  │             │    │ 8KB         │    │ 100:3389    │         │
│  └─────────────┘    └─────────────┘    └─────────────┘         │
│           │                 │                 │                │
│           ▼                 ▼                 ▼                │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │                    Metrics & Monitoring                     │ │
│  │                                                             │ │
│  │ • Bytes transferred                                         │ │
│  │ • Active connections                                        │ │
│  │ • Connection duration                                       │ │
│  │ • Error rates                                               │ │
│  └─────────────────────────────────────────────────────────────┘ │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### TCP Proxy Features

#### 1. **Bidirectional Data Transfer**
```go
// Двунаправленная передача данных
func (p *TCPProxy) handleConnection(clientConn, serverConn net.Conn) {
    // Client → Server
    go p.transferData(clientConn, serverConn, "client_to_server")
    // Server → Client
    go p.transferData(serverConn, clientConn, "server_to_client")
}
```

#### 2. **Buffer Management**
```go
// Управление буферами
type BufferManager struct {
    BufferSize    int           `json:"buffer_size"`
    MaxBuffers    int           `json:"max_buffers"`
    BufferPool    chan []byte   `json:"buffer_pool"`
}
```

#### 3. **Connection Pooling**
```go
// Пул соединений
type ConnectionPool struct {
    MaxConnections int                    `json:"max_connections"`
    Connections    map[string]*Connection `json:"connections"`
    Mutex          sync.RWMutex           `json:"mutex"`
}
```

## 🏢 Multi-tenancy Архитектура

### Модель Multi-tenancy

```
┌─────────────────────────────────────────────────────────────────┐
│                    Multi-tenancy Model                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   Tenant A      │  │   Tenant B      │  │   Tenant C      │ │
│  │                 │  │                 │  │                 │ │
│  │ • Tunnels: 5/10 │  │ • Tunnels: 2/20 │  │ • Tunnels: 0/5  │ │
│  │ • Connections:  │  │ • Connections:  │  │ • Connections:  │ │
│  │   50/100        │  │   150/200       │  │   0/50          │ │
│  │ • Bandwidth:    │  │ • Bandwidth:    │  │ • Bandwidth:    │ │
│  │   50/100 Mbps   │  │   300/500 Mbps  │  │   0/100 Mbps    │ │
│  │ • IP Range:     │  │ • IP Range:     │  │ • IP Range:     │ │
│  │   192.168.1.0/24│  │   192.168.2.0/24│  │   192.168.3.0/24│ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
├─────────────────────────────────────────────────────────────────┤
│                    Shared Infrastructure                       │
│  • CPU, Memory, Network                                        │
│  • TLS Certificates                                            │
│  • Monitoring & Metrics                                        │
│  • Rate Limiting                                               │
└─────────────────────────────────────────────────────────────────┘
```

### Multi-tenancy Features

#### 1. **Resource Isolation**
```yaml
tenants:
  tenant_001:
    limits:
      max_tunnels: 20
      max_connections_per_tunnel: 200
      max_bandwidth_mbps: 500
      max_data_transfer_gb: 5000
      rate_limit_per_second: 2000
```

#### 2. **IP Filtering**
```yaml
tenants:
  tenant_001:
    allowed_ips:
      - "192.168.1.0/24"
      - "10.0.0.0/8"
    blocked_ips:
      - "192.168.1.100"
      - "10.0.0.50"
```

#### 3. **Activity Tracking**
```go
type TenantActivity struct {
    TenantID      string    `json:"tenant_id"`
    LastActivity  time.Time `json:"last_activity"`
    ActiveTunnels int       `json:"active_tunnels"`
    TotalTraffic  int64     `json:"total_traffic"`
}
```

## 🔐 mTLS Архитектура

### Модель взаимной аутентификации

```
┌─────────────────────────────────────────────────────────────────┐
│                    mTLS Authentication Model                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐         │
│  │   Desktop   │    │   Relay     │    │   Server    │         │
│  │   Client    │    │   Server    │    │   Client    │         │
│  │             │    │             │    │             │         │
│  │ • Client    │◄──►│ • Server    │◄──►│ • Client    │         │
│  │   Cert      │    │   Cert      │    │   Cert      │         │
│  │ • Private   │    │ • Private   │    │ • Private   │         │
│  │   Key       │    │   Key       │    │   Key       │         │
│  │ • CA Chain  │    │ • CA Chain  │    │ • CA Chain  │         │
│  └─────────────┘    └─────────────┘    └─────────────┘         │
│           │                 │                 │                │
│           ▼                 ▼                 ▼                │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │                    Certificate Authority                   │ │
│  │                                                             │ │
│  │ • Root CA (self-signed)                                    │ │
│  │ • Intermediate CA (for clients)                            │ │
│  │ • Certificate Revocation List (CRL)                       │ │
│  │ • Online Certificate Status Protocol (OCSP)               │ │
│  └─────────────────────────────────────────────────────────────┘ │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Процесс mTLS handshake

#### 1. **Client Hello**
```json
{
  "client_certificate": "-----BEGIN CERTIFICATE-----...",
  "client_key_exchange": "RSA/ECDHE",
  "cipher_suites": ["TLS_AES_256_GCM_SHA384", "TLS_CHACHA20_POLY1305_SHA256"],
  "supported_groups": ["P-256", "P-384", "X25519"]
}
```

#### 2. **Server Hello**
```json
{
  "server_certificate": "-----BEGIN CERTIFICATE-----...",
  "selected_cipher": "TLS_AES_256_GCM_SHA384",
  "selected_group": "P-256",
  "session_id": "unique-session-id"
}
```

#### 3. **Certificate Verification**
```go
// Проверка сертификатов
func verifyCertificate(cert *x509.Certificate, roots *x509.CertPool) error {
    // Проверка подписи
    if err := cert.CheckSignatureFrom(roots); err != nil {
        return fmt.Errorf("certificate signature verification failed: %v", err)
    }
    
    // Проверка срока действия
    if time.Now().Before(cert.NotBefore) || time.Now().After(cert.NotAfter) {
        return fmt.Errorf("certificate expired or not yet valid")
    }
    
    // Проверка отзыва (CRL/OCSP)
    if err := checkRevocation(cert); err != nil {
        return fmt.Errorf("certificate revoked: %v", err)
    }
    
    return nil
}
```

### Конфигурация mTLS

```json
{
  "tls": {
    "enabled": true,
    "mode": "mutual",
    "cert_file": "/etc/relay/certs/server.crt",
    "key_file": "/etc/relay/certs/server.key",
    "ca_file": "/etc/relay/certs/ca.crt",
    "client_ca_file": "/etc/relay/certs/client-ca.crt",
    "cipher_suites": [
      "TLS_AES_256_GCM_SHA384",
      "TLS_CHACHA20_POLY1305_SHA256"
    ],
    "min_version": "TLS 1.3",
    "session_tickets": true,
    "session_cache_size": 32768
  }
}
```

## 🌐 Edge Networks Architecture

### Распределенная сеть edge-узлов

```
┌─────────────────────────────────────────────────────────────────┐
│                    Global Edge Network                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐         │
│  │   Europe    │    │   Americas  │    │   Asia      │         │
│  │   Region    │    │   Region    │    │   Region    │         │
│  │             │    │             │    │             │         │
│  │ ┌─────────┐ │    │ ┌─────────┐ │    │ ┌─────────┐ │         │
│  │ │ Edge-1  │ │    │ │ Edge-1  │ │    │ │ Edge-1  │ │         │
│  │ │ DE      │ │    │ │ US-East │ │    │ │ JP      │ │         │
│  │ └─────────┘ │    │ └─────────┘ │    │ └─────────┘ │         │
│  │             │    │             │    │             │         │
│  │ ┌─────────┐ │    │ ┌─────────┐ │    │ ┌─────────┐ │         │
│  │ │ Edge-2  │ │    │ │ Edge-2  │ │    │ │ Edge-2  │ │         │
│  │ │ UK      │ │    │ │ US-West │ │    │ │ SG      │ │         │
│  │ └─────────┘ │    │ └─────────┘ │    │ └─────────┘ │         │
│  └─────────────┘    └─────────────┘    └─────────────┘         │
│           │                 │                 │                │
│           ▼                 ▼                 ▼                │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │                    Global Load Balancer                    │ │
│  │                                                             │ │
│  │ • Geo-based routing                                        │ │
│  │ • Health checks                                            │ │
│  │ • Failover logic                                           │ │
│  │ • Traffic optimization                                     │ │
│  └─────────────────────────────────────────────────────────────┘ │
│           │                                                     │
│           ▼                                                     │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │                    Core Infrastructure                     │ │
│  │                                                             │ │
│  │ • Relay Server Cluster                                     │ │
│  │ • Keycloak HA Setup                                        │ │
│  │ • Django API Cluster                                       │ │
│  │ • Database Replication                                     │ │
│  └─────────────────────────────────────────────────────────────┘ │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Edge Node Configuration

```yaml
# edge-node-config.yaml
edge_nodes:
  - name: "edge-eu-de-1"
    region: "europe"
    country: "DE"
    location: "Frankfurt"
    capacity: "10k_connections"
    upstream:
      - "relay-core-1:8080"
      - "relay-core-2:8080"
    health_check:
      interval: "30s"
      timeout: "5s"
      unhealthy_threshold: 3
      healthy_threshold: 2
    load_balancing:
      algorithm: "least_connections"
      sticky_sessions: true
      session_timeout: "1h"
    security:
      rate_limit: "1000_rps"
      ddos_protection: true
      ip_whitelist: ["10.0.0.0/8", "172.16.0.0/12"]
    monitoring:
      metrics_endpoint: ":8081"
      log_level: "INFO"
      alerting: true
```

### Geo-routing и failover

```go
// Geo-based routing logic
func selectEdgeNode(clientIP net.IP, userAgent string) (*EdgeNode, error) {
    // Определение региона клиента
    region := geoip.Lookup(clientIP)
    
    // Выбор edge-узла в регионе
    edgeNodes := getEdgeNodesInRegion(region)
    
    // Проверка доступности
    availableNodes := filterHealthyNodes(edgeNodes)
    
    if len(availableNodes) == 0 {
        // Failover на ближайший регион
        return selectFailoverNode(region)
    }
    
    // Выбор по алгоритму балансировки
    return selectByLoadBalancingAlgorithm(availableNodes)
}
```

## 🚀 Production Load Handling

### High Availability (HA) Setup

#### Relay Server Cluster
```yaml
# relay-cluster.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: relay-server
spec:
  replicas: 3
  serviceName: relay-service
  selector:
    matchLabels:
      app: relay-server
  template:
    metadata:
      labels:
        app: relay-server
    spec:
      containers:
      - name: relay
        image: cloudbridge/relay:latest
        ports:
        - containerPort: 8080
        - containerPort: 8081  # metrics
        env:
        - name: RELAY_CLUSTER_MODE
          value: "true"
        - name: RELAY_CLUSTER_NODES
          value: "relay-0,relay-1,relay-2"
        - name: RELAY_ELECTION_TIMEOUT
          value: "5s"
        volumeMounts:
        - name: relay-config
          mountPath: /etc/relay
        - name: relay-certs
          mountPath: /etc/relay/certs
        livenessProbe:
          httpGet:
            path: /health
            port: 8081
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8081
          initialDelaySeconds: 5
          periodSeconds: 5
```

#### Auto Scaling Configuration
```yaml
# relay-hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: relay-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: StatefulSet
    name: relay-server
  minReplicas: 3
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
  - type: Object
    object:
      metric:
        name: active_connections
      describedObject:
        apiVersion: v1
        kind: Service
        name: relay-service
      target:
        type: AverageValue
        averageValue: 1000
  behavior:
    scaleUp:
      stabilizationWindowSeconds: 60
      policies:
      - type: Percent
        value: 100
        periodSeconds: 15
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
      - type: Percent
        value: 10
        periodSeconds: 60
```

### Load Testing и Capacity Planning

#### Load Testing Scenarios
```python
# load_test_scenarios.py
import asyncio
import aiohttp
import time

class LoadTest:
    def __init__(self, relay_url, num_clients=1000):
        self.relay_url = relay_url
        self.num_clients = num_clients
        self.results = []
    
    async def test_connection_establishment(self):
        """Тест установки соединений"""
        start_time = time.time()
        tasks = []
        
        for i in range(self.num_clients):
            task = asyncio.create_task(self.establish_connection(i))
            tasks.append(task)
        
        await asyncio.gather(*tasks)
        duration = time.time() - start_time
        
        return {
            "test": "connection_establishment",
            "clients": self.num_clients,
            "duration": duration,
            "rate": self.num_clients / duration
        }
    
    async def test_tunnel_throughput(self):
        """Тест пропускной способности туннелей"""
        # Реализация теста пропускной способности
        pass
    
    async def test_concurrent_tunnels(self):
        """Тест одновременных туннелей"""
        # Реализация теста одновременных туннелей
        pass

# Сценарии тестирования
scenarios = [
    {"name": "baseline", "clients": 100, "duration": 300},
    {"name": "peak_load", "clients": 1000, "duration": 600},
    {"name": "stress_test", "clients": 5000, "duration": 900},
    {"name": "endurance", "clients": 500, "duration": 3600}
]
```

#### Performance Metrics
```promql
# Метрики производительности
# Пропускная способность туннелей
rate(relay_tunnel_bytes_transferred_total[5m])

# Задержка соединений
histogram_quantile(0.95, rate(relay_connection_duration_seconds_bucket[5m]))

# Использование ресурсов
relay_cpu_usage_percent
relay_memory_usage_bytes
relay_network_io_bytes_total

# Ошибки и отказы
rate(relay_connection_errors_total[5m])
rate(relay_authentication_failures_total[5m])

# Автоскейлинг метрики
relay_active_connections
relay_pending_connections
relay_tunnel_creation_duration_seconds
```

### Misuse Protection

#### Многоуровневая защита
```go
// Misuse Protection Implementation
type MisuseProtector struct {
    rateLimiters map[string]*rate.Limiter
    ipBlacklist  *sync.Map
    userScores   *sync.Map
    rules        []MisuseRule
}

type MisuseRule struct {
    Name        string
    Condition   func(*Connection) bool
    Action      MisuseAction
    Threshold   int
    Window      time.Duration
}

type MisuseAction struct {
    Type        string // "block", "throttle", "alert", "disable"
    Duration    time.Duration
    Message     string
}

// Правила защиты от злоупотреблений
var misuseRules = []MisuseRule{
    {
        Name: "rapid_connections",
        Condition: func(conn *Connection) bool {
            return conn.rate > 100 // > 100 connections per minute
        },
        Action: MisuseAction{
            Type:     "throttle",
            Duration: 5 * time.Minute,
            Message:  "Too many rapid connections",
        },
        Threshold: 5,
        Window:    1 * time.Minute,
    },
    {
        Name: "failed_auth_attempts",
        Condition: func(conn *Connection) bool {
            return conn.authFailures > 10
        },
        Action: MisuseAction{
            Type:     "block",
            Duration: 30 * time.Minute,
            Message:  "Too many failed authentication attempts",
        },
        Threshold: 3,
        Window:    10 * time.Minute,
    },
    {
        Name: "suspicious_patterns",
        Condition: func(conn *Connection) bool {
            return detectSuspiciousPattern(conn)
        },
        Action: MisuseAction{
            Type:     "alert",
            Duration: 0,
            Message:  "Suspicious connection pattern detected",
        },
        Threshold: 1,
        Window:    1 * time.Hour,
    },
}
```

#### DDoS Protection
```go
// DDoS Protection
type DDoSProtector struct {
    ipCounters   *sync.Map
    globalLimiter *rate.Limiter
    whitelist    map[string]bool
    blacklist    map[string]bool
}

func (d *DDoSProtector) CheckConnection(conn *Connection) error {
    // Проверка белого списка
    if d.whitelist[conn.IP] {
        return nil
    }
    
    // Проверка черного списка
    if d.blacklist[conn.IP] {
        return errors.New("IP is blacklisted")
    }
    
    // Глобальное ограничение
    if !d.globalLimiter.Allow() {
        return errors.New("global rate limit exceeded")
    }
    
    // Проверка по IP
    counter := d.getIPCounter(conn.IP)
    if counter.Count() > 1000 { // 1000 connections per minute per IP
        d.blacklist[conn.IP] = true
        return errors.New("IP rate limit exceeded")
    }
    
    return nil
}
```

## Слои архитектуры

### A. Control Plane (управляющий уровень)
* Keycloak
* Django API
* Биллинг
* Админ-панель
* Панель клиента
* API Gateway

### B. Data Plane (трафик)
* Relay Server
* Server Client
* Desktop Client
* TCP-туннели

### C. Observability
* Prometheus
* Grafana
* AlertManager
* Логирование, tracing

### D. Security Intelligence
* SIEM
* ML Engine
* Реакция на угрозы (rules + AI)
* Интеграция с внешними DLP/SOC

## Преимущества архитектуры

* Единый источник доверия — все токены управляются Keycloak
* Безопасность на основе открытых стандартов (OAuth2, OIDC)
* Масштабируемость — можно подключить любые внешние IdP
* Полная изоляция данных между клиентами
* Аудит — все действия логируются в Keycloak и Django
* Поведенческий анализ и ML-детекция угроз
* Интеграция с SIEM и внешними системами безопасности

## Компоненты

### Relay Server

* Принимает соединения от Server Client и Desktop Client
* Валидирует токены через Keycloak (JWKS)
* Создаёт и управляет туннелями
* Проксирует TCP трафик
* Поддержка TLS, rate limiting, мониторинга

### Server Client

* Запрашивает токен через Django API
* Подключается к Relay Server
* Передаёт информацию о туннелях (remote_host, remote_port)
* Проксирует трафик от Relay Server во внутреннюю сеть

### Desktop Client

* Выполняет вход через Keycloak (OIDC flow)
* Получает access token
* Подключается к Relay Server и использует туннели по правам доступа

### Django API

* Предоставляет REST endpoint для генерации токена Server Client
* Делегирует генерацию токена Keycloak
* Хранит информацию о серверах, конфигурации, мониторинге

### Keycloak

* Централизованная система авторизации
* Поддерживает локальных пользователей, LDAP, внешние IdP
* Подписывает JWT токены
* Выдаёт access и refresh токены
* JWKS используется Relay Server для проверки подлинности токенов

## Потоки

### Аутентификация Server Client

```
Админ создает сервер (Django UI)
→ Server Client запрашивает токен (POST /api/token)
→ Django → Keycloak (client credentials grant)
→ Получение токена
→ Server Client → Relay Server (с токеном)
→ Relay Server → JWKS → Валидация
→ Туннель создан
```

### Аутентификация Desktop Client

```
Пользователь логинится через Keycloak (OIDC)
→ Получает access token
→ Подключается к Relay Server
→ Relay Server проверяет подпись токена
→ Доступ разрешён
```

## Пример токена

```json
{
  "sub": "user@example.com",
  "aud": "relay-client",
  "exp": 1710000000,
  "iat": 1709990000,
  "client_id": "server-123",
  "tunnels": [
    {
      "id": "ldap-main",
      "remote_host": "192.168.0.10",
      "remote_port": 389
    }
  ]
}
```

## Протокол соединения с Relay Server

### 1. Установка соединения и приветствие

```json
{
  "type": "hello",
  "version": "1.0"
}
```

Ответ:

```json
{
  "type": "hello_response",
  "version": "1.0"
}
```

### 2. Аутентификация токеном

```json
{
  "type": "auth",
  "token": "<JWT от Keycloak>",
  "role": "server|desktop"
}
```

Ответ:

```json
{
  "type": "auth_response",
  "status": "ok",
  "client_id": "server-123"
}
```

### 3. Информация о туннелях (только для Server Client)

```json
{
  "type": "tunnel_info",
  "tunnels": [
    {
      "id": "ldap-main",
      "remote_host": "192.168.0.10",
      "remote_port": 389
    }
  ]
}
```

Ответ:

```json
{
  "type": "tunnel_response",
  "status": "ok",
  "external_port": 3890
}
```

### 4. Heartbeat

```json
{ "type": "heartbeat" }
```

Ответ:

```json
{ "type": "heartbeat_response" }
```

## Конфигурация валидации токенов (Relay Server)

```json
{
  "keycloak": {
    "enabled": true,
    "jwks_url": "https://keycloak.example.com/realms/cloudbridge/protocol/openid-connect/certs",
    "issuer": "https://keycloak.example.com/realms/cloudbridge",
    "audience": "relay-client"
  }
}
```

## Минимальные требования токена

* Подпись: RS256
* Claims:
  * `sub` — идентификатор пользователя
  * `client_id` — ID сервера или пользователя
  * `exp`, `iat` — срок действия
  * `aud` — audience (`relay-client`)
  * `tunnels` — список разрешённых туннелей (для Server Client)

## Пример запроса токена (Server Client → Django → Keycloak)

```bash
curl -X POST https://api.example.com/api/relay/token \
  -H "Authorization: Bearer <admin_token>" \
  -d '{"server_id": "123"}'
```

Django получает токен от Keycloak через client_credentials flow и возвращает его Server Client.

## Пример входа (Desktop Client → Keycloak)

```bash
curl -X POST https://keycloak.example.com/realms/cloudbridge/protocol/openid-connect/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password" \
  -d "client_id=relay-client" \
  -d "username=user@example.com" \
  -d "password=password"
```

## Защита и контроль

* Все соединения шифруются (TLS 1.3)
* Все токены валидируются через Keycloak JWKS
* Включён rate limiting и защита от DDoS
* Поддержка мультидоменных реалмов и LDAP-федерации
* Изоляция туннелей по `client_id` из токена

## Мониторинг

* Prometheus метрики по туннелям, токенам, соединениям
* Grafana дашборды
* Алерты через AlertManager
* Аудит доступа — по логам Django и Keycloak

## 🚀 Основные возможности

- ✅ **Туннелирование TCP-соединений** с поддержкой TLS/mTLS
- ✅ **JWT аутентификация** через Keycloak (OAuth 2.0 / OpenID Connect)
- ✅ **AD/LDAP интеграция** через Keycloak User Federation
- ✅ **Мультитенантность** с изоляцией данных
- ✅ **Балансировка нагрузки** (round-robin)
- ✅ **Ограничение скорости запросов** (rate limiting)
- ✅ **Расширенный мониторинг** через Prometheus + Grafana
- ✅ **Алерты и уведомления** через AlertManager
- ✅ **Структурированное логирование** (JSON) с ротацией
- ✅ **Автоматическое переподключение**
- ✅ **Graceful shutdown**
- ✅ **Health checks**
- ✅ **Полное тестовое покрытие** (unit, integration, E2E)
- ✅ **SIEM интеграция** с поведенческим анализом
- ✅ **ML-детекция угроз** и аномалий
- ✅ **Open API Gateway** для внешних интеграций
- ✅ **Webhooks** для событий безопасности
- ✅ **mTLS (Mutual TLS)** - двусторонняя аутентификация
- ✅ **Edge Networks** - глобальная распределенная сеть
- ✅ **High Availability (HA)** - кластеризация и отказоустойчивость
- ✅ **Auto Scaling** - автоматическое масштабирование
- ✅ **Misuse Protection** - защита от злоупотреблений
- ✅ **DDoS Protection** - защита от распределенных атак
- ✅ **Load Testing** - тестирование производительности
- ✅ **Capacity Planning** - планирование ресурсов

## 🔒 Безопасность

### Обязательные действия после установки

1. **Измените пароли по умолчанию**:
   ```bash
   # Используйте скрипт генерации секретов
   ./scripts/generate_secrets.sh --env .env
   
   # Или измените вручную в файле .env:
   KEYCLOAK_ADMIN_PASSWORD=your-secure-password
   GRAFANA_ADMIN_PASSWORD=your-secure-password
   ```

2. **Настройте HTTPS**:
   - Замените самоподписанные сертификаты на реальные
   - Обновите `RELAY_TLS_CERT_FILE` и `RELAY_TLS_KEY_FILE`

3. **Настройте firewall**:
   ```bash
   sudo ufw allow 22/tcp    # SSH
   sudo ufw allow 8080/tcp  # Keycloak Admin Console
   sudo ufw allow 8082/tcp  # Relay API
   sudo ufw allow 8081/tcp  # Metrics
   sudo ufw enable
   ```

4. **Используйте менеджер секретов в продакшене**:
   - HashiCorp Vault
   - AWS Secrets Manager
   - Azure Key Vault
   - Google Secret Manager

### Генерация секретов

```bash
# Генерация всех секретов
./scripts/generate_secrets.sh --all

# Пароль указанной длины
./scripts/generate_secrets.sh --password 20

# .env файл с секретами
./scripts/generate_secrets.sh --env .env.production
```

### Соответствие российским требованиям

- ✅ **152-ФЗ**: Защита персональных данных
- ✅ **187-ФЗ**: Безопасность критической инфраструктуры
- ✅ **ГОСТ Р 34.10-2012**: Алгоритмы подписи
- ✅ **ГОСТ Р 34.11-2012**: Алгоритмы хеширования

## 📊 Мониторинг

### Prometheus метрики
- Активные соединения и туннели
- Производительность и задержки
- Ошибки и отказы
- Использование ресурсов
- Keycloak метрики (валидация токенов, JWKS кэш)
- Бизнес-метрики (активные пользователи, доход)

### Grafana дашборды
- Готовые дашборды для мониторинга
- Автоматическая настройка через provisioning
- Визуализация метрик в реальном времени

### Алерты
- Настраиваемые алерты Prometheus
- Уведомления через Slack, Email
- Escalation policies
- Автоматическое разрешение

## 📝 Логирование

### Структурированное логирование
- **Формат:** JSON (по умолчанию), Console
- **Уровни:** DEBUG, INFO, WARN, ERROR
- **Ротация:** Автоматическая с настройкой размера и количества файлов
- **Поля:** timestamp, level, message, user_id, request_id, connection_id, duration, error

### Специализированные методы
- Логирование соединений и аутентификации
- Логирование валидации токенов и проверки разрешений
- Логирование производительности и безопасности
- Логирование бизнес-событий

## 🧪 Тестирование

### Типы тестов
- **Unit тесты** - покрытие всех компонентов
- **Integration тесты** - тестирование взаимодействия
- **E2E тесты** - полный цикл работы
- **Performance тесты** - нагрузочное тестирование

### Покрытие
- Более 80% покрытия кода
- Автоматизированные тесты в CI/CD
- Тестирование безопасности

## 🛠️ Установка

### Требования

- Go 1.18+
- Docker 20.10+ (опционально)
- Systemd (для установки как сервис)

### Быстрый старт

### Для московского региона

```bash
# Клонирование репозитория
git clone https://github.com/your-org/cloudbridge-relay-installer.git
cd cloudbridge-relay-installer

# Генерация безопасных секретов
./scripts/generate_secrets.sh --env .env

# Редактирование конфигурации (при необходимости)
nano .env

# Запуск сервисов
docker-compose -f docker-compose.keycloak.yml up -d
```

Этот скрипт автоматически:
- Создаст `.env` файл с настройками для московского региона
- Настроит часовой пояс Europe/Moscow
- Установит русскую локаль
- Создаст необходимые директории
- Настроит сертификаты
- Запустит все сервисы с помощью Docker Compose
- Проверит соответствие российским требованиям

### Стандартная установка

```bash
# Клонирование репозитория
git clone https://github.com/your-org/cloudbridge-relay-installer.git
cd cloudbridge-relay-installer

# Копирование примера конфигурации
cp env.example .env

# Редактирование конфигурации
nano .env

# Запуск сервисов
docker-compose -f docker-compose.keycloak.yml up -d
```

## ⚙️ Конфигурация

### Основные параметры

```json
{
  "server": {
    "port": 8080,
    "max_tunnels": 1000,
    "timeout": "30s",
    "ping_interval": "30s",
    "cleanup_interval": "1m"
  },
  "keycloak": {
    "enabled": true,
    "jwks_url": "https://keycloak.example.com/realms/cloudbridge/protocol/openid-connect/certs",
    "issuer": "https://keycloak.example.com/realms/cloudbridge",
    "audience": "relay-client",
    "token_validation": {
      "clock_skew": "30s"
    }
  },
  "security": {
    "tls_enabled": true,
    "cert_file": "/etc/relay/certs/fullchain.pem",
    "key_file": "/etc/relay/certs/privkey.pem",
    "rate_limit": {
      "enabled": true,
      "requests_per_second": 100,
      "burst": 200
    }
  },
  "metrics": {
    "enabled": true,
    "port": ":8081"
  },
  "logging": {
    "level": "INFO",
    "format": "json",
    "output": "/var/log/relay/relay.log"
  }
}
```

**Важно:** Все поля типа `time.Duration` должны быть указаны в виде строк с суффиксами времени (например, "30s", "5m", "1h").

### Переменные окружения

```bash
# Опциональные
RELAY_LOG_LEVEL=info
RELAY_TLS_ENABLED=true

# Keycloak
KEYCLOAK_SERVER_URL=https://keycloak.example.com
KEYCLOAK_REALM=cloudbridge
KEYCLOAK_CLIENT_ID=relay-client
KEYCLOAK_CLIENT_SECRET=your-client-secret

# Moscow Region Configuration
DEPLOYMENT_REGION=moscow
TZ=Europe/Moscow
LANG=ru_RU.UTF-8
LC_ALL=ru_RU.UTF-8

# Russia-Specific Configuration
COMPLIANCE_152_FZ=true
COMPLIANCE_187_FZ=true
USE_RUSSIAN_CRYPTO=true
```

## 📖 Документация

### Основная документация
- [API Reference](docs/api_reference.md) - Полная документация API
- [Deployment Guide](docs/deployment_guide.md) - Руководство по развертыванию
- [Client Integration](docs/client_integration.md) - Интеграция клиентов
- [Debugging](docs/debugging.md) - Отладка и диагностика
- [Metrics](docs/metrics.md) - Архитектура и использование метрик
- [Keycloak Integration](docs/keycloak_integration.md) - Интеграция с Keycloak

### Примеры использования
- [Python Client](examples/python_client.py)
- [JavaScript Client](examples/javascript_client.js)
- [Go Client](examples/go_client.go)

## 🔧 Разработка

### Команды Makefile

```bash
# Сборка
make build

# Тестирование
make test-all

# Линтинг
make lint

# Форматирование
make fmt

# Docker
make docker-build
make docker-run

# Мониторинг
make setup-monitoring

# Безопасность
make security-audit
```

### Структура проекта

```
cloudbridge-relay-installer/
├── cmd/                    # Исполняемые файлы
├── internal/               # Внутренние пакеты
│   ├── auth/              # Аутентификация
│   ├── keycloak/          # Интеграция с Keycloak
│   ├── metrics/           # Метрики
│   └── tunnel/            # Управление туннелями
├── config/                # Конфигурация
├── scripts/               # Скрипты автоматизации
├── tests/                 # Тесты
└── docs/                  # Документация
```

## 🚨 Устранение неполадок

### Частые проблемы

#### 1. Ошибка парсинга конфигурации
```
Failed to load configuration: failed to parse config file: json: cannot unmarshal string into Go struct field .server.timeout of type time.Duration
```
**Решение:** Убедитесь, что все поля `time.Duration` указаны в виде строк с суффиксами времени.

#### 2. Порт уже занят
```
Failed to start server: failed to create listener: listen tcp :8080: bind: address already in use
```
**Решение:** Остановите процесс, использующий порт 8080, или измените порт в конфигурации.

#### 3. Ошибка Docker контейнера
```
docker: Error response from daemon: failed to create task for container: failed to create shim task: OCI runtime create failed: runc create failed: unable to start container process: error during container init: exec: "--config": executable file not found in $PATH
```
**Решение:** Укажите правильную команду запуска: `cloudbridge-relay ./relay --config config/config.json`

### Диагностика

```bash
# Проверка статуса
sudo systemctl status relay

# Просмотр логов
tail -f /var/log/relay/relay.log

# Проверка метрик
curl http://localhost:8081/metrics

# Проверка конфигурации
./relay --config config/config.json --help
```

## 📈 Мониторинг

### Основные метрики

```promql
# Активные соединения
relay_active_connections

# Процент успешных валидаций токенов
rate(relay_keycloak_token_validations_total{result="success"}[5m]) / 
rate(relay_keycloak_token_validations_total[5m]) * 100

# Время соединения (95-й перцентиль)
histogram_quantile(0.95, rate(relay_connection_duration_seconds_bucket[5m]))
```

### Алерты

- Высокий процент неудачных валидаций токенов
- Медленные операции JWKS
- Низкая производительность кэша
- Высокое использование памяти

## 🤝 Поддержка

### Получение помощи

1. Проверьте [документацию](docs/)
2. Изучите [часто задаваемые вопросы](docs/debugging.md)
3. Создайте issue в репозитории

### Сбор информации для поддержки

```bash
# Логи
sudo journalctl -u relay > relay-logs.txt

# Метрики
curl http://localhost:8081/metrics > metrics.txt

# Конфигурация
cat /etc/relay/config.json > config.txt

# Системная информация
uname -a > system-info.txt
```

## 📄 Лицензия

Этот проект лицензирован под MIT License - см. файл [LICENSE](LICENSE) для деталей.

## 🏗️ Архитектура

CloudBridge Relay построен с использованием современных практик разработки:

- **Модульная архитектура** - четкое разделение ответственности
- **Dependency Injection** - легкое тестирование и конфигурация
- **Graceful Shutdown** - корректное завершение работы
- **Health Checks** - мониторинг состояния сервиса
- **Structured Logging** - удобный анализ логов
- **Metrics First** - детальный мониторинг производительности

## 🚨 Управление

### Основные команды

```bash
# Запуск сервисов
docker-compose -f docker-compose.keycloak.yml up -d

# Остановка сервисов
docker-compose -f docker-compose.keycloak.yml down

# Просмотр логов
docker-compose -f docker-compose.keycloak.yml logs -f

# Перезапуск сервиса
docker-compose -f docker-compose.keycloak.yml restart relay

# Обновление сервисов
docker-compose -f docker-compose.keycloak.yml pull
docker-compose -f docker-compose.keycloak.yml up -d
```

### Резервное копирование

```bash
# Создание резервной копии
./scripts/backup_keycloak.sh

# Восстановление из резервной копии
./scripts/restore_keycloak.sh backup_file.tar.gz
```

## 🆚 Сравнение с Cloudflare Argo Tunnel

| Функция | Cloudflare Argo Tunnel | CloudBridge Relay |
|---------|------------------------|-------------------|
| **Контроль инфраструктуры** | ❌ Зависимость от Cloudflare | ✅ Полный контроль |
| **Данные** | ❌ Проходят через Cloudflare | ✅ Ваши серверы |
| **Стоимость** | ❌ Платная подписка | ✅ Собственная инфраструктура |
| **Кастомизация** | ❌ Ограниченная | ✅ Полная |
| **Мультитенантность** | ❌ Базовая | ✅ Продвинутая |
| **AD/LDAP интеграция** | ❌ Отсутствует | ✅ Полная поддержка |
| **Мониторинг** | ❌ Cloudflare Analytics | ✅ Prometheus + Grafana |
| **Соответствие требованиям** | ❌ Зависит от юрисдикции | ✅ Полный контроль |
| **SIEM интеграция** | ❌ Отсутствует | ✅ Полная поддержка |
| **ML-детекция угроз** | ❌ Отсутствует | ✅ Встроенная |
| **mTLS аутентификация** | ❌ Отсутствует | ✅ Полная поддержка |
| **Edge Networks** | ❌ Только Cloudflare PoP | ✅ Собственная сеть |
| **High Availability** | ❌ Зависит от Cloudflare | ✅ Полный контроль |
| **Auto Scaling** | ❌ Ограниченное | ✅ Полная автоматизация |
| **Misuse Protection** | ❌ Базовая | ✅ Продвинутая защита |
| **DDoS Protection** | ✅ Встроенная | ✅ Настраиваемая |
| **Load Testing** | ❌ Отсутствует | ✅ Встроенные инструменты |

## Особенности

- 🔐 Интеграция с Keycloak для централизованной аутентификации
- 🚀 Высокая производительность и масштабируемость
- 📊 Встроенный мониторинг с Prometheus и Grafana
- 🔒 Поддержка TLS/SSL шифрования
- 📝 Подробное логирование
- 🛡️ Rate limiting и защита от DDoS
- 🇷🇺 Поддержка российских требований (152-ФЗ, 187-ФЗ)
- 🏢 **Мультитенантность с изоляцией компаний**
- 🔄 **Туннелирование любых TCP сервисов**
- 👥 **Интеграция с Active Directory/LDAP**
- 🎯 **Полный контроль над инфраструктурой**
- 🧠 **ML-детекция угроз и поведенческий анализ**
- 📊 **SIEM интеграция с внешними системами**
- 🔌 **Open API Gateway для интеграций**
- 🔐 **mTLS (Mutual TLS) аутентификация**
- 🌐 **Глобальная сеть edge-узлов**
- 🚀 **High Availability и Auto Scaling**
- 🛡️ **Продвинутая защита от злоупотреблений**
- 📈 **Load Testing и Capacity Planning**