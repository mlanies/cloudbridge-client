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

> **Внимание:** Не храните токены и пароли в открытом виде. Используйте переменные окружения и защищенные хранилища для секретов.

## Особенности

- Автоматическая установка и настройка
- Поддержка systemd
- Автоматическое обновление
- Настраиваемые пути установки
- Подробное логирование
- Проверка системных требований
- Кроссплатформенная поддержка (Linux, Windows, macOS)
- Поддержка российских ОС (Astra Linux, Alt Linux, ROSA Linux, RedOS) через бинарник для Linux

## Требования

- Linux с systemd (для systemd-интеграции)
- curl
- root права (для установки как системный сервис)

## Быстрая установка

### Linux/macOS/Российские ОС

```bash
# Скачать последнюю версию
curl -L https://github.com/2gc-dev/cloudbridge-client/releases/latest/download/cloudbridge-client-linux-amd64 -o cloudbridge-client
chmod +x cloudbridge-client
sudo mv cloudbridge-client /usr/local/bin/

# Или установить через Go
go install github.com/2gc-dev/cloudbridge-client/cmd/cloudbridge-client@latest
```

> **Примечание:** Для российских дистрибутивов (Astra Linux, Alt Linux, ROSA Linux, RedOS и др.) используйте бинарник cloudbridge-client-linux-amd64.

### Windows

```powershell
# Скачать последнюю версию
Invoke-WebRequest -Uri "https://github.com/2gc-dev/cloudbridge-client/releases/latest/download/cloudbridge-client-windows-amd64.exe" -OutFile "cloudbridge-client.exe"
```

## Сборка из исходников

```bash
# Клонировать репозиторий
git clone https://github.com/2gc-dev/cloudbridge-client.git
cd cloudbridge-client

# Собрать для текущей платформы
make build

# Собрать для всех платформ
make build-all
```

## Usage

### Basic Usage

```bash
Usage: cloudbridge-client [options]

Options:
  -h, --help                 Показать справку
  -v, --version              Показать версию
  -c, --config PATH         Путь к конфигурационному файлу
```

## Примеры

Запуск с пользовательским конфигурационным файлом:
```bash
cloudbridge-client --config /path/to/config.yaml
```

## После установки

### Linux (systemd)

Клиент будет установлен и запущен как системный сервис. Вы можете управлять им с помощью systemd:

```bash
# Проверить статус
systemctl status cloudbridge-client

# Посмотреть логи
journalctl -u cloudbridge-client -f

# Перезапустить сервис
systemctl restart cloudbridge-client
```

### Windows

Для Windows рекомендуется использовать NSSM для установки как службы:

```powershell
# Установить через NSSM
nssm install CloudBridgeClient "C:\path\to\cloudbridge-client.exe"
nssm set CloudBridgeClient AppParameters "--config C:\path\to\config.yaml"
nssm start CloudBridgeClient
```

## Конфигурация

Конфигурационный файл находится в `/etc/cloudbridge-client/config.yaml` (Linux) или в указанном месте. После изменения конфигурации перезапустите сервис.

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

## License

This project is licensed under the MIT License - see the LICENSE file for details. 