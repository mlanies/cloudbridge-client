#!/bin/bash

# =============================================================================
# 2GC Installer Script
# =============================================================================
#
# Этот скрипт предназначен для установки и настройки компонентов 2GC:
# - Cloudflare tunnel (Zero Trust)
# - CloudBridge Client
#
# Поддерживаемые операционные системы:
# - macOS (10.15 и новее)
# - Linux (systemd-based дистрибутивы)
#
# Требования:
# - Root права (sudo)
# - curl
# - На macOS: Homebrew (для установки cloudflared)
#
# Использование:
#   1. Сделайте скрипт исполняемым:
#      chmod +x installer.sh
#
#   2. Запустите с правами root:
#      sudo ./installer.sh
#
#   3. Следуйте инструкциям на экране:
#      - Выберите компонент для установки
#      - Введите ваш 2GC Token
#
# Особенности:
# - Автоматическое определение ОС и архитектуры
# - Проверка существующих установок
# - Возможность удаления старых версий
# - Цветной вывод для лучшей читаемости
# - Подробное логирование процесса установки
#
# Пути установки:
# - Cloudflared: /usr/local/bin/cloudflared
# - CloudBridge Client: /usr/local/bin/cloudbridge-client
#
# Автор: 2GC
# Версия: 1.0.0
# Лицензия: MIT
# =============================================================================

# Цвета для вывода
CYAN='\033[0;36m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Проверка root прав
if [ "$EUID" -ne 0 ]; then 
    echo -e "${RED}Этот скрипт должен быть запущен с правами root (sudo)${NC}"
    exit 1
fi

# Красивый заголовок
echo -e "\n"
echo -e "${CYAN}╔════════════════════════════════════╗"
echo -e "║          2GC INSTALLER             ║"
echo -e "║      ╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌           ║"
echo -e "║ Ваша безопасность — наш приоритет  ║"
echo -e "╚════════════════════════════════════╝${NC}"
echo -e "\n"
echo -e "${YELLOW}    Безопасность — это не роскошь,"
echo -e "            это стандарт.${NC}"
echo -e "\n"
echo -e " © 2GC, 2025 | https://2gc.ru"
echo -e "\n"

# Определение ОС
if [[ "$OSTYPE" == "darwin"* ]]; then
    OS="macos"
elif [[ -f /etc/os-release ]]; then
    . /etc/os-release
    OS=$ID
else
    echo -e "${RED}Не удалось определить операционную систему${NC}"
    exit 1
fi

# Функция для установки Cloudflared
install_cloudflared() {
    echo -e "\n${CYAN}[Cloudflared] Проверка существующей установки...${NC}"
    
    if command -v cloudflared &> /dev/null; then
        echo -e "${YELLOW}[Cloudflared] Уже установлен.${NC}"
        read -p "Удалить старую версию? (y/n) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            if [ "$OS" = "macos" ]; then
                brew uninstall cloudflared
            else
                systemctl stop cloudflared
                systemctl disable cloudflared
                rm -f /usr/local/bin/cloudflared
            fi
        else
            echo "Отмена установки Cloudflared."
            return
        fi
    fi

    echo -e "${CYAN}[Cloudflared] Установка...${NC}"
    if [ "$OS" = "macos" ]; then
        brew install cloudflared
    else
        curl -L https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64 -o /usr/local/bin/cloudflared
        chmod +x /usr/local/bin/cloudflared
    fi

    echo -e "${CYAN}[Cloudflared] Регистрация токена...${NC}"
    cloudflared service install "$TOKEN"
    echo -e "${GREEN}[Cloudflared] Установка и регистрация завершены!${NC}"
}

# Функция для установки CloudBridge Client
install_cloudbridge() {
    echo -e "\n${CYAN}[CloudBridge Client] Проверка существующей установки...${NC}"
    
    if [ -f "/usr/local/bin/cloudbridge-client" ]; then
        echo -e "${YELLOW}[CloudBridge Client] Уже установлен.${NC}"
        read -p "Удалить старую версию? (y/n) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            cloudbridge-client service uninstall
            rm -f /usr/local/bin/cloudbridge-client
        else
            echo "Отмена установки CloudBridge Client."
            return
        fi
    fi

    echo -e "${CYAN}[CloudBridge Client] Скачивание дистрибутива...${NC}"
    if [ "$OS" = "macos" ]; then
        curl -L https://github.com/mlanies/cloudbridge-client/releases/latest/download/cloudbridge-client-darwin-amd64 -o /usr/local/bin/cloudbridge-client
    else
        curl -L https://github.com/mlanies/cloudbridge-client/releases/latest/download/cloudbridge-client-linux-amd64 -o /usr/local/bin/cloudbridge-client
    fi

    if [ ! -f "/usr/local/bin/cloudbridge-client" ]; then
        echo -e "${RED}Файл не найден после скачивания.${NC}"
        return
    fi

    chmod +x /usr/local/bin/cloudbridge-client
    echo -e "${CYAN}[CloudBridge Client] Установка и регистрация...${NC}"
    cloudbridge-client service install "$TOKEN"
    echo -e "${GREEN}[CloudBridge Client] Установка и регистрация завершены!${NC}"
}

# Меню выбора
echo "Выберите, что хотите зарегистрировать:"
echo " 1 - Cloudflare tunnel (Zero Trust)"
echo " 2 - CloudBridge Client"
read -p "Введите цифру (1/2): " choice

# Запрос токена
read -p "Введите ваш 2GC Token: " TOKEN
if [ -z "$TOKEN" ]; then
    echo -e "${RED}Токен не введён. Выход.${NC}"
    exit 1
fi

# Установка выбранного компонента
case $choice in
    1)
        install_cloudflared
        ;;
    2)
        install_cloudbridge
        ;;
    *)
        echo -e "${RED}Некорректный выбор. Выход.${NC}"
        exit 1
        ;;
esac

echo -e "\n${GREEN}Готово!${NC}" 