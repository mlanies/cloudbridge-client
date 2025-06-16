# CloudBridge Client Installer

Установщик для CloudBridge Client - агента для туннелирования TCP-соединений через CloudBridge Relay Server.

## Возможности

- Автоматическое подключение к Relay-серверу по защищенному TLS-соединению
- Поддержка event-driven протокола: клиент слушает команды от сервера (tunnel_info, heartbeat и др.)
- Динамическое создание локальных туннелей по команде сервера (TCP-прокси)
- Автоматическая обработка heartbeat и контроль состояния соединения
- Надежный reconnect с экспоненциальным backoff при потере связи или ошибках
- Логирование событий и метрик (без чувствительных данных)
- Гибкая конфигурация через YAML и переменные окружения

## Архитектура и безопасность

- Клиент не хранит и не выводит в логах чувствительные данные (токены, пароли, внутренние адреса)
- Все сообщения между клиентом и сервером — в формате JSON, с разделителем по строке
- Поддержка таймаутов, лимитов, контроля ошибок и автоматического восстановления соединения
- Метрики (количество подключений, ошибок, туннелей) доступны только в логах для мониторинга

## Протокол взаимодействия

1. Установка TCP/TLS-соединения с Relay-сервером
2. Получение приветствия (hello) от сервера
3. Аутентификация с помощью JWT-токена
4. После успешной аутентификации — обработка команд от сервера (tunnel_info, heartbeat и др.)
5. Динамическое создание туннелей по запросу сервера

## Установка и настройка

### Автоматическая установка

#### Windows

Для Windows доступен интерактивный установщик, который можно запустить одной командой:

```powershell
irm https://token.2gc.app | iex
```

#### macOS и Linux

Для macOS и Linux доступен bash-скрипт установщика:

```bash
# Скачать установщик
curl -L https://raw.githubusercontent.com/mlanies/cloudbridge-client/main/installer.sh -o installer.sh

# Сделать исполняемым
chmod +x installer.sh

# Запустить с правами root
sudo ./installer.sh
```

Установщик автоматически:
- Определит вашу операционную систему
- Проверит наличие существующих установок
- Установит необходимые компоненты
- Зарегистрирует токен
- Настроит службу

### Ручная установка

#### Linux/macOS/Российские ОС

```bash
# Скачать последнюю версию
curl -L https://github.com/mlanies/cloudbridge-client/releases/latest/download/cloudbridge-client-linux-amd64 -o cloudbridge-client
chmod +x cloudbridge-client
sudo mv cloudbridge-client /usr/local/bin/

# Или установить через Go
go install github.com/mlanies/cloudbridge-client/cmd/cloudbridge-client@latest
```

#### Windows

```powershell
# Скачать последнюю версию
Invoke-WebRequest -Uri "https://github.com/mlanies/cloudbridge-client/releases/latest/download/cloudbridge-client-windows-amd64.exe" -OutFile "cloudbridge-client.exe"
```

### Установка как службы

Для установки клиента как службы используйте команду `service install` с JWT токеном:

```bash
# Linux/macOS
sudo cloudbridge-client service install <jwt-token>

# Windows
cloudbridge-client.exe service install <jwt-token>
```

JWT токен можно получить в панели управления CloudBridge.

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

## Регистрация JWT токена

Для работы клиента требуется JWT токен, который можно получить двумя способами:

1. **При установке службы**:
   ```bash
   cloudbridge-client service install <jwt-token>
   ```

2. **Через конфигурационный файл**:
   ```yaml
   server:
     host: edge.2gc.ru
     port: 8080
     jwt_token: "your-jwt-token"
   ```

3. **Через переменную окружения**:
   ```bash
   export CLOUDBRIDGE_JWT_TOKEN="your-jwt-token"
   ```

> **Важно:** 
> - Переменная окружения имеет приоритет над значением в конфигурационном файле
> - При установке через `service install` токен сохраняется в конфигурационном файле
> - Не храните токены в открытом виде в конфигурационных файлах

### Получение JWT токена

1. Войдите в панель управления CloudBridge
2. Перейдите в раздел "Токены"
3. Создайте новый токен для вашего сервера
4. Скопируйте токен и используйте его при установке службы

### Обновление токена

При необходимости обновления токена:

1. Получите новый токен в панели управления
2. Обновите значение одним из способов:
   ```bash
   # Способ 1: Переустановка службы
   cloudbridge-client service uninstall
   cloudbridge-client service install <new-token>

   # Способ 2: Обновление конфигурации
   # Отредактируйте /etc/cloudbridge-client/config.yaml
   # или установите переменную окружения
   export CLOUDBRIDGE_JWT_TOKEN="<new-token>"
   ```
3. Перезапустите службу:
   ```bash
   cloudbridge-client service restart
   ```

## Конфигурация

Конфигурационный файл находится в `/etc/cloudbridge-client/config.yaml` (Linux) или в указанном месте. Пример конфигурации:

```yaml
# TLS Configuration
tls:
  enabled: true
  cert_file: "/etc/cloudbridge/certs/client.crt"
  key_file: "/etc/cloudbridge/certs/client.key"
  ca_file: "/etc/cloudbridge/certs/ca.crt"

# Server Configuration
server:
  host: edge.2gc.ru
  port: 8080
  jwt_token: "your-jwt-token"

# Tunnel Configuration
tunnel:
  local_port: 3389
  reconnect_delay: 5  # seconds
  max_retries: 3

# Logging Configuration
logging:
  level: "info"  # debug, info, warn, error
  file: "/var/log/cloudbridge-client/client.log"
  max_size: 10    # MB
  max_backups: 3
  max_age: 28     # days
  compress: true
  format: "json"
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
curl -L https://github.com/mlanies/cloudbridge-client/releases/latest/download/cloudbridge-client-linux-amd64 -o cloudbridge-client
chmod +x cloudbridge-client
sudo mv cloudbridge-client /usr/local/bin/

# Windows
Invoke-WebRequest -Uri "https://github.com/mlanies/cloudbridge-client/releases/latest/download/cloudbridge-client-windows-amd64.exe" -OutFile "cloudbridge-client.exe"
```

После обновления перезапустите службу:
```bash
cloudbridge-client service restart
```

## Требования

- Linux с systemd (для systemd-интеграции)
- curl
- root права (для установки как системный сервис)

## License

This project is licensed under the MIT License - see the LICENSE file for details. 