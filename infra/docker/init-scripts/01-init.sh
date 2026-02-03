#!/bin/bash
# Скрипт инициализации базы данных PostgreSQL

set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- Включить UUID расширение
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
    
    -- Создать схемы для микросервисов
    CREATE SCHEMA IF NOT EXISTS auth;
    CREATE SCHEMA IF NOT EXISTS schedule;
    CREATE SCHEMA IF NOT EXISTS ticket;
    CREATE SCHEMA IF NOT EXISTS fiscal;
    CREATE SCHEMA IF NOT EXISTS payment;
    CREATE SCHEMA IF NOT EXISTS board;
    CREATE SCHEMA IF NOT EXISTS geo;
    CREATE SCHEMA IF NOT EXISTS notify;
    CREATE SCHEMA IF NOT EXISTS audit;
    CREATE SCHEMA IF NOT EXISTS document;
    
    -- Создать пользователей для сервисов (опционально)
    CREATE USER vokzal_auth WITH PASSWORD 'auth_pass_2026';
    CREATE USER vokzal_ticket WITH PASSWORD 'ticket_pass_2026';
    CREATE USER vokzal_schedule WITH PASSWORD 'schedule_pass_2026';
    
    -- Выдать права
    GRANT ALL PRIVILEGES ON SCHEMA auth TO vokzal_auth;
    GRANT ALL PRIVILEGES ON SCHEMA ticket TO vokzal_ticket;
    GRANT ALL PRIVILEGES ON SCHEMA schedule TO vokzal_schedule;
    
    -- Общие привилегии
    GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO admin;
    GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO admin;
    
    -- Вывести информацию
    SELECT 'Вокзал.ТЕХ database initialized successfully!' AS message;
EOSQL
