# CloudBridge Relay Client

Кроссплатформенный Go-клиент для CloudBridge Relay с поддержкой TLS 1.3, JWT-аутентификации, multi-tenancy и комплексной обработкой ошибок. Клиент реализует полную спецификацию протокола согласно техническому заданию v2.0.

## Возможности

- **TLS 1.3**: Принудительное использование TLS 1.3 с безопасными cipher suites
- **JWT-аутентификация**: Полная валидация JWT-токенов с поддержкой HMAC и RSA
- **Multi-tenancy**: Поддержка tenant_id в JWT токенах и изоляция ресурсов
- **Интеграция с Keycloak**: Опциональная поддержка OpenID Connect
- **Кроссплатформенность**: Windows, Linux, macOS (x86_64, ARM64)
- **Rate Limiting**: Встроенное ограничение скорости с экспоненциальным backoff
- **Heartbeat**: Автоматический мониторинг состояния соединения
- **Управление туннелями**: Полный жизненный цикл туннелей с буферизацией
- **Prometheus метрики**: Комплексная система мониторинга производительности
- **Оптимизация производительности**: Автоматическая оптимизация для высокой пропускной способности
- **Обработка ошибок**: Комплексная обработка ошибок и логика повторных попыток
- **Конфигурация**: Гибкая YAML-конфигурация с поддержкой переменных окружения

## Новые возможности v2.0

### Multi-tenancy
- Поддержка `tenant_id` в JWT токенах
- Изоляция ресурсов по тенантам
- Передача `tenant_id` в сообщениях туннеля

### Улучшенный TCP Proxy
- Управление буферами с пулом соединений
- Статистика передачи данных
- Оптимизация производительности

### Система метрик
- Prometheus интеграция
- Метрики по тенантам и туннелям
- Мониторинг производительности в реальном времени

### Оптимизация производительности
- Автоматическая оптимизация runtime
- Поддержка высоко- и низко-латентных режимов
- Управление сборщиком мусора

## Протокол взаимодействия

Клиент реализует полный протокол CloudBridge Relay:

1. **Hello/Hello Response**: Согласование версии протокола
2. **Auth/Auth Response**: JWT-аутентификация
3. **Tunnel Info/Tunnel Response**: Создание и управление туннелями
4. **Heartbeat/Heartbeat Response**: Мониторинг состояния соединения
5. **Error Messages**: Стандартизированная обработка ошибок

## Установка

### Автоматическая установка

#### Windows

```powershell
irm https://token.2gc.app | iex
```

#### macOS и Linux

```bash
# Скачать установщик
curl -L https://raw.githubusercontent.com/2gc-dev/cloudbridge-client/main/installer.sh -o installer.sh

# Сделать исполняемым
chmod +x installer.sh

# Запустить с правами root
sudo ./installer.sh
```

### Ручная установка

#### Через Go Install

```bash
go install github.com/2gc-dev/cloudbridge-client/cmd/cloudbridge-client@latest
```

#### Готовые бинарники

Скачайте подходящий бинарник для вашей платформы со [страницы релизов](https://github.com/2gc-dev/cloudbridge-client/releases).

#### Сборка из исходников

```bash
git clone https://github.com/2gc-dev/cloudbridge-client.git
cd cloudbridge-client
go build -o cloudbridge-client ./cmd/cloudbridge-client
```

### Настройка окружения разработки

Для разработки требуется Go 1.21.13 или выше:

```bash
# Установка Go 1.21.13 (если не установлен)
wget https://go.dev/dl/go1.21.13.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.21.13.linux-amd64.tar.gz

# Настройка PATH
export PATH=/usr/local/go/bin:$PATH
export PATH=$HOME/go/bin:$PATH

# Или используйте готовый скрипт
source setup-go.sh

# Установка golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.57.2

# Проверка установки
go version
golangci-lint version
```

### Запуск тестов и линтера

```bash
# Запуск тестов
go test -v ./...

# Запуск линтера
golangci-lint run

# Сборка проекта
go build -v -o cloudbridge-client ./cmd/cloudbridge-client
```

### Система релизов

Проект использует автоматическую систему релизов через GitHub Actions:

#### Автоматические релизы
- При каждом push в `main` создается автоматический релиз с версией `{VERSION}-{COMMIT_HASH}`
- Все тесты и линтинг должны пройти успешно
- Создаются бинарные файлы для всех платформ (Windows, Linux, macOS)

#### Ручные релизы
Для создания релиза с тегом:

```bash
# Обновить версию (patch, minor, major)
./scripts/bump-version.sh patch

# Или создать тег вручную
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3
```

#### Версионирование
- Файл `VERSION` содержит текущую версию проекта
- Автоматические релизы: `{VERSION}-{COMMIT_HASH}` (например: `0.1.0-a1b2c3d`)
- Ручные релизы: `v{VERSION}` (например: `v1.2.3`)

#### Платформы
Каждый релиз включает бинарные файлы для:
- **Windows**: x64, ARM64
- **Linux**: x64, ARM64  
- **macOS**: x64, ARM64

## Быстрый старт

### Базовое использование

```bash
cloudbridge-client --token "your-jwt-token"
```

Подключится к серверу по умолчанию (edge.2gc.ru:8080) с включенным TLS.

### С конфигурационным файлом

```bash
cloudbridge-client --config config.yaml --token "your-jwt-token"
```

### Пользовательский туннель

```bash
cloudbridge-client \
  --token "your-jwt-token" \
  --tunnel-id "my-tunnel" \
  --local-port 3389 \
  --remote-host "192.168.1.100" \
  --remote-port 3389
```

## Конфигурация

Клиент поддерживает конфигурацию через YAML-файлы и переменные окружения.

### Конфигурационный файл (config.yaml)

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

### Переменные окружения

Все опции конфигурации можно установить через переменные окружения с префиксом `CLOUDBRIDGE_`:

```bash
export CLOUDBRIDGE_RELAY_HOST="edge.2gc.ru"
export CLOUDBRIDGE_RELAY_PORT="8080"
export CLOUDBRIDGE_AUTH_SECRET="your-jwt-secret"
```

### Параметры командной строки

- `--config, -c`: Путь к конфигурационному файлу
- `--token, -t`: JWT-токен для аутентификации (обязательно)
- `--tunnel-id, -i`: ID туннеля (по умолчанию: tunnel_001)
- `--local-port, -l`: Локальный порт для привязки (по умолчанию: 3389)
- `--remote-host, -r`: Удаленный хост (по умолчанию: 192.168.1.100)
- `--remote-port, -p`: Удаленный порт (по умолчанию: 3389)
- `--verbose, -v`: Включить подробное логирование

## Установка как службы

### Установка службы

```bash
# Linux/macOS
sudo cloudbridge-client service install <jwt-token>

# Windows
cloudbridge-client.exe service install <jwt-token>
```

### Управление службой

```bash
# Проверка статуса
cloudbridge-client service status

# Запуск службы
cloudbridge-client service start

# Остановка службы
cloudbridge-client service stop

# Перезапуск службы
cloudbridge-client service restart

# Удаление службы
cloudbridge-client service uninstall
```

## Безопасность

### TLS 1.3

- Принудительное использование TLS 1.3
- Только безопасные cipher suites:
  - `TLS_AES_256_GCM_SHA384`
  - `TLS_CHACHA20_POLY1305_SHA256`
  - `TLS_AES_128_GCM_SHA256`
- Валидация сертификатов
- Поддержка SNI

### JWT-аутентификация

- Поддержка HMAC-SHA256
- Валидация RSA-подписи
- Проверка срока действия токена
- Извлечение subject для rate limiting

### Интеграция с Keycloak

- Поддержка OpenID Connect
- Автоматическое получение JWKS
- Валидация токенов
- Контроль доступа на основе ролей

## Обработка ошибок

Клиент обрабатывает все стандартные ошибки relay:

- `invalid_token`: Неверный или истекший JWT-токен
- `rate_limit_exceeded`: Ограничение скорости с экспоненциальным backoff
- `connection_limit_reached`: Превышен лимит соединений
- `server_unavailable`: Недоступность сервера с повторными попытками
- `invalid_tunnel_info`: Неверная конфигурация туннеля
- `unknown_message_type`: Ошибки протокола

## Rate Limiting

Встроенное ограничение скорости с настраиваемыми параметрами:

- Ограничение по пользователю (на основе JWT subject)
- Стратегия экспоненциального backoff
- Настраиваемое максимальное количество попыток
- Ограничения максимального backoff

## Heartbeat

Автоматический мониторинг состояния соединения:

- Настраиваемый интервал heartbeat (по умолчанию: 30с)
- Обнаружение и обработка сбоев
- Автоматическое переподключение при сбоях
- Статистика heartbeat

## Поддерживаемые платформы

Протестировано и поддерживается на:

- **Windows**: x86_64, ARM64
- **Linux**: x86_64, ARM64
- **macOS**: x86_64, ARM64

## Разработка

### Сборка для нескольких платформ

```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o cloudbridge-client.exe ./cmd/cloudbridge-client

# Linux
GOOS=linux GOARCH=amd64 go build -o cloudbridge-client ./cmd/cloudbridge-client

# macOS
GOOS=darwin GOARCH=amd64 go build -o cloudbridge-client ./cmd/cloudbridge-client
```

### Запуск тестов

```bash
go test ./...
```

### Структура кода

```
pkg/
├── auth/          # Управление аутентификацией
├── config/        # Обработка конфигурации
├── errors/        # Обработка ошибок и логика повторных попыток
├── heartbeat/     # Управление heartbeat
├── interfaces/    # Интерфейсы для модульности
├── relay/         # Основной relay-клиент
├── service/       # Управление службами
├── tunnel/        # Управление туннелями
└── types/         # Типы данных

cmd/
└── cloudbridge-client/  # Основное приложение

docs/
├── API.md              # API документация
├── ARCHITECTURE.md     # Архитектура системы
├── DEPLOYMENT.md       # Инструкции по развертыванию
├── job.md              # Техническое задание
├── PERFORMANCE.md      # Производительность
├── README.md           # Основная документация
├── SECURITY.md         # Безопасность
├── TESTING.md          # Тестирование
└── TROUBLESHOOTING.md  # Устранение неполадок
```

## Мониторинг

### Логи

```bash
# Linux (systemd)
journalctl -u cloudbridge-client -f

# Windows
# Логи доступны через Event Viewer

# macOS
tail -f /var/log/cloudbridge-client/client.log
```

### Метрики

Метрики доступны по адресу `http://localhost:9090/metrics` в формате Prometheus.

## Обновление

Для обновления до последней версии:

```bash
# Linux/macOS
curl -L https://github.com/2gc-dev/cloudbridge-client/releases/latest/download/cloudbridge-client-linux-amd64 -o cloudbridge-client
chmod +x cloudbridge-client
sudo mv cloudbridge-client /usr/local/bin/

# Windows
Invoke-WebRequest -Uri "https://github.com/2gc-dev/cloudbridge-client/releases/latest/download/cloudbridge-client-windows-amd64.exe" -OutFile "cloudbridge-client.exe"
```

После обновления перезапустите службу:
```bash
cloudbridge-client service restart
```

## Участие в разработке

1. Форкните репозиторий
2. Создайте ветку для новой функции (`git checkout -b feature/amazing-feature`)
3. Зафиксируйте изменения (`git commit -m 'Add some amazing feature'`)
4. Отправьте в ветку (`git push origin feature/amazing-feature`)
5. Откройте Pull Request

## Лицензия

Этот проект лицензирован под MIT License - см. файл LICENSE для деталей.

## Поддержка

Для поддержки и вопросов:

- Создайте issue на GitHub
- Изучите документацию
- Просмотрите примеры конфигурации

## Changelog

### v1.1.1
- Исправлены тесты туннелей для соответствия архитектуре
- Удалены устаревшие и неработающие тесты
- Приведение проекта в соответствие с техническим заданием
- Улучшена документация и структура кода

### v1.0.0
- Первоначальный релиз
- Поддержка TLS 1.3
- JWT-аутентификация
- Кроссплатформенная поддержка
- Комплексная обработка ошибок
- Rate limiting и логика повторных попыток
- Механизм heartbeat
- Управление туннелями

## Требования

- Linux с systemd (для systemd-интеграции)
- curl
- root права (для установки как системный сервис)
