# CloudBridge Client Installer

Установщик для CloudBridge Client - агента для туннелирования TCP-соединений через CloudBridge Relay Server.

## Особенности

- Автоматическая установка и настройка
- Поддержка systemd
- Автоматическое обновление
- Настраиваемые пути установки
- Подробное логирование
- Проверка системных требований
- Кроссплатформенная поддержка (Linux, Windows, macOS)
- Поддержка российских ОС (Astra Linux, Alt Linux, ROSA Linux, RedOS)

## Требования

- Linux с systemd (для systemd-интеграции)
- curl
- root права (для установки как системный сервис)

## Быстрая установка

### Linux/macOS

```bash
# Скачать последнюю версию
curl -L https://github.com/2gc-dev/cloudbridge-client/releases/latest/download/cloudbridge-client-linux-amd64 -o cloudbridge-client
chmod +x cloudbridge-client
sudo mv cloudbridge-client /usr/local/bin/

# Или установить через Go
go install github.com/2gc-dev/cloudbridge-client/cmd/cloudbridge-client@latest
```

### Российские ОС

#### Astra Linux
```bash
curl -L https://github.com/2gc-dev/cloudbridge-client/releases/latest/download/cloudbridge-client-astra-amd64 -o cloudbridge-client
chmod +x cloudbridge-client
sudo mv cloudbridge-client /usr/local/bin/
```

#### Alt Linux
```bash
curl -L https://github.com/2gc-dev/cloudbridge-client/releases/latest/download/cloudbridge-client-alt-amd64 -o cloudbridge-client
chmod +x cloudbridge-client
sudo mv cloudbridge-client /usr/local/bin/
```

#### ROSA Linux
```bash
curl -L https://github.com/2gc-dev/cloudbridge-client/releases/latest/download/cloudbridge-client-rosa-amd64 -o cloudbridge-client
chmod +x cloudbridge-client
sudo mv cloudbridge-client /usr/local/bin/
```

#### RedOS
```bash
curl -L https://github.com/2gc-dev/cloudbridge-client/releases/latest/download/cloudbridge-client-redos-amd64 -o cloudbridge-client
chmod +x cloudbridge-client
sudo mv cloudbridge-client /usr/local/bin/
```

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

# Собрать только для российских ОС
make build-russian
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