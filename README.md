# CloudBridge Client Installer

Установщик для CloudBridge Client - агента для туннелирования TCP-соединений через CloudBridge Relay Server.

## Особенности

- Автоматическая установка и настройка
- Поддержка systemd
- Автоматическое обновление
- Настраиваемые пути установки
- Подробное логирование
- Проверка системных требований

## Требования

- Linux с systemd
- curl
- root права

## Быстрая установка

```bash
curl -sSL https://raw.githubusercontent.com/twogc/cloudbridge-client-installer/main/install.sh | sudo bash
```

## Ручная установка

1. Скачайте установщик:
```bash
curl -O https://raw.githubusercontent.com/twogc/cloudbridge-client-installer/main/install.sh
```

2. Сделайте его исполняемым:
```bash
chmod +x install.sh
```

3. Запустите установку:
```bash
sudo ./install.sh
```

## Опции установки

```bash
Usage: install.sh [options]

Options:
  -h, --help                 Показать справку
  -v, --version              Показать версию
  -d, --debug                Включить отладочный вывод
  --client-version VERSION   Установить конкретную версию (по умолчанию: latest)
  --install-dir DIR          Директория установки (по умолчанию: /usr/local/bin)
  --config-dir DIR           Директория конфигурации (по умолчанию: /etc/cloudbridge-client)
  --data-dir DIR             Директория данных (по умолчанию: /var/lib/cloudbridge-client)
  --log-dir DIR              Директория логов (по умолчанию: /var/log/cloudbridge-client)
```

## Примеры

Установка конкретной версии:
```bash
sudo ./install.sh --client-version 1.0.0
```

Установка с отладочным выводом:
```bash
sudo ./install.sh --debug
```

Установка в пользовательскую директорию:
```bash
sudo ./install.sh --install-dir /opt/cloudbridge-client
```

## После установки

Клиент будет установлен и запущен как системный сервис. Вы можете управлять им с помощью systemd:

```bash
# Проверить статус
systemctl status cloudbridge-client

# Посмотреть логи
journalctl -u cloudbridge-client -f

# Перезапустить сервис
systemctl restart cloudbridge-client
```

## Конфигурация

Конфигурационный файл находится в `/etc/cloudbridge-client/config.yaml`. После изменения конфигурации перезапустите сервис:

```bash
sudo systemctl restart cloudbridge-client
```

## Обновление

Для обновления до последней версии:

```bash
sudo ./install.sh
```

Для обновления до конкретной версии:

```bash
sudo ./install.sh --client-version 1.1.0
```

## Удаление

Для удаления клиента:

```bash
sudo systemctl stop cloudbridge-client
sudo systemctl disable cloudbridge-client
sudo rm /etc/systemd/system/cloudbridge-client.service
sudo rm -rf /etc/cloudbridge-client
sudo rm -rf /var/lib/cloudbridge-client
sudo rm -rf /var/log/cloudbridge-client
sudo rm /usr/local/bin/cloudbridge-client
sudo userdel cloudbridge-client
sudo systemctl daemon-reload
```

## Лицензия

MIT License 