#!/bin/bash
set -e

# Скрипт для безопасной настройки администратора Вокзал.ТЕХ
# Использование: ./scripts/setup-admin.sh

echo "==================================="
echo "Вокзал.ТЕХ - Настройка администратора"
echo "==================================="
echo ""

# Проверка переменных окружения
if [ -z "$DB_HOST" ]; then
    echo "Введите хост базы данных (по умолчанию: localhost):"
    read -r DB_HOST
    DB_HOST=${DB_HOST:-localhost}
fi

if [ -z "$DB_PORT" ]; then
    echo "Введите порт базы данных (по умолчанию: 5432):"
    read -r DB_PORT
    DB_PORT=${DB_PORT:-5432}
fi

if [ -z "$DB_NAME" ]; then
    echo "Введите имя базы данных (по умолчанию: vokzal):"
    read -r DB_NAME
    DB_NAME=${DB_NAME:-vokzal}
fi

if [ -z "$DB_USER" ]; then
    echo "Введите пользователя базы данных (по умолчанию: postgres):"
    read -r DB_USER
    DB_USER=${DB_USER:-postgres}
fi

if [ -z "$DB_PASSWORD" ]; then
    echo "Введите пароль базы данных:"
    read -rs DB_PASSWORD
    echo ""
fi

echo ""
echo "==================================="
echo "Настройка пароля администратора"
echo "==================================="
echo ""
echo "Требования к паролю:"
echo "- Минимум 12 символов"
echo "- Содержит заглавные и строчные буквы"
echo "- Содержит цифры"
echo "- Содержит специальные символы"
echo ""

# Функция проверки надежности пароля
validate_password() {
    local password="$1"
    local length=${#password}
    
    if [ $length -lt 12 ]; then
        echo "Ошибка: Пароль должен быть не менее 12 символов"
        return 1
    fi
    
    if ! echo "$password" | grep -q '[A-Z]'; then
        echo "Ошибка: Пароль должен содержать заглавные буквы"
        return 1
    fi
    
    if ! echo "$password" | grep -q '[a-z]'; then
        echo "Ошибка: Пароль должен содержать строчные буквы"
        return 1
    fi
    
    if ! echo "$password" | grep -q '[0-9]'; then
        echo "Ошибка: Пароль должен содержать цифры"
        return 1
    fi
    
    if ! echo "$password" | grep -q '[!@#$%^&*()_+=-]'; then
        echo "Ошибка: Пароль должен содержать специальные символы (!@#$%^&*()_+=- и т.д.)"
        return 1
    fi
    
    return 0
}

# Ввод и подтверждение пароля
while true; do
    echo "Введите новый пароль администратора:"
    read -rs NEW_PASSWORD
    echo ""
    
    if ! validate_password "$NEW_PASSWORD"; then
        echo ""
        continue
    fi
    
    echo "Подтвердите пароль:"
    read -rs NEW_PASSWORD_CONFIRM
    echo ""
    
    if [ "$NEW_PASSWORD" != "$NEW_PASSWORD_CONFIRM" ]; then
        echo "Ошибка: Пароли не совпадают"
        echo ""
        continue
    fi
    
    break
done

echo ""
echo "Генерация безопасного хэша пароля..."

# Генерация bcrypt хэша с помощью Go
# Передаём пароль через переменную окружения для предотвращения инъекции команд
PASSWORD_HASH=$(docker run --rm -e NEW_PASSWORD="$NEW_PASSWORD" golang:1.22-alpine sh -c '
go install golang.org/x/crypto/bcrypt/cmd/bcrypt@latest && \
echo -n "$NEW_PASSWORD" | /root/go/bin/bcrypt
' 2>/dev/null | tail -n 1)

if [ -z "$PASSWORD_HASH" ]; then
    echo "Ошибка: Не удалось сгенерировать хэш пароля"
    echo "Попробуйте установить Go и выполнить:"
    echo "  go install golang.org/x/crypto/bcrypt/cmd/bcrypt@latest"
    echo "  echo -n 'ваш_пароль' | bcrypt"
    exit 1
fi

echo "Хэш сгенерирован успешно"
echo ""
echo "Обновление администратора в базе данных..."

# Обновление пароля в БД с использованием переменных psql для безопасности
export PGPASSWORD="$DB_PASSWORD"
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" \
  -v password_hash="$PASSWORD_HASH" <<'EOF'
UPDATE users 
SET password_hash = :'password_hash', 
    is_active = true,
    updated_at = NOW()
WHERE id = '00000000-0000-0000-0000-000000000001';
EOF

# Если psql завершился успешно (set -e остановит скрипт при ошибке), выводим сообщение
echo ""
echo "✅ Администратор успешно настроен!"
echo ""
echo "Данные для входа:"
echo "  Логин: admin"
echo "  Пароль: <введенный вами пароль>"
echo ""
echo "URL админ-панели: http://localhost:3001"
