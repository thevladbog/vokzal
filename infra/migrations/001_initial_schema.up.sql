-- Migration: 001_initial_schema
-- Description: Создание основных таблиц для Вокзал.ТЕХ

-- Включить расширения
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm"; -- для полнотекстового поиска

-- ====================================
-- Таблица: stations (Автовокзалы)
-- ====================================
CREATE TABLE stations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(10) UNIQUE NOT NULL,
    address TEXT,
    timezone VARCHAR(50) NOT NULL DEFAULT 'Europe/Moscow',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_stations_code ON stations(code);
COMMENT ON TABLE stations IS 'Автовокзалы';

-- ====================================
-- Таблица: buses (Автобусы)
-- ====================================
CREATE TABLE buses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    plate_number VARCHAR(12) UNIQUE NOT NULL,
    model VARCHAR(50) NOT NULL,
    capacity INTEGER NOT NULL CHECK (capacity > 0),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'maintenance', 'out_of_service')),
    station_id UUID NOT NULL REFERENCES stations(id) ON DELETE RESTRICT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_buses_station ON buses(station_id);
CREATE INDEX idx_buses_status ON buses(status);
COMMENT ON TABLE buses IS 'Автобусы';

-- ====================================
-- Таблица: drivers (Водители)
-- ====================================
CREATE TABLE drivers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    full_name VARCHAR(100) NOT NULL,
    license_number VARCHAR(20) UNIQUE NOT NULL,
    experience_years INTEGER CHECK (experience_years >= 0),
    phone VARCHAR(15),
    station_id UUID NOT NULL REFERENCES stations(id) ON DELETE RESTRICT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_drivers_station ON drivers(station_id);
CREATE INDEX idx_drivers_license ON drivers(license_number);
COMMENT ON TABLE drivers IS 'Водители';

-- ====================================
-- Таблица: routes (Маршруты)
-- ====================================
CREATE TABLE routes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    stops JSONB NOT NULL,
    distance_km DECIMAL(8,2),
    duration_min INTEGER,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_routes_active ON routes(is_active);
CREATE INDEX idx_routes_stops ON routes USING GIN (stops);
COMMENT ON TABLE routes IS 'Маршруты';
COMMENT ON COLUMN routes.stops IS 'JSON: [{ "station_id": "...", "order": 1, "arrival_offset_min": 0 }]';

-- ====================================
-- Таблица: schedules (Расписание)
-- ====================================
CREATE TABLE schedules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    route_id UUID NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
    departure_time TIME NOT NULL,
    days_of_week JSONB NOT NULL,
    platform VARCHAR(10),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_schedules_route ON schedules(route_id);
CREATE INDEX idx_schedules_active ON schedules(is_active);
COMMENT ON TABLE schedules IS 'Расписание рейсов';
COMMENT ON COLUMN schedules.days_of_week IS 'JSON: [1,3,5] - дни недели (1=пн, 7=вс)';

-- ====================================
-- Таблица: trips (Рейсы)
-- ====================================
CREATE TABLE trips (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    schedule_id UUID NOT NULL REFERENCES schedules(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'delayed', 'cancelled', 'departed', 'arrived')),
    delay_minutes INTEGER DEFAULT 0,
    platform VARCHAR(10),
    departure_actual TIMESTAMP,
    arrival_actual TIMESTAMP,
    bus_id UUID REFERENCES buses(id) ON DELETE SET NULL,
    driver_id UUID REFERENCES drivers(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_schedule_date UNIQUE (schedule_id, date)
);

CREATE INDEX idx_trips_schedule_date ON trips(schedule_id, date);
CREATE INDEX idx_trips_status ON trips(status);
CREATE INDEX idx_trips_date ON trips(date);
COMMENT ON TABLE trips IS 'Рейсы (экземпляры расписания)';

-- ====================================
-- Таблица: seats (Места в автобусе)
-- ====================================
CREATE TABLE seats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bus_id UUID NOT NULL REFERENCES buses(id) ON DELETE CASCADE,
    number INTEGER NOT NULL CHECK (number > 0),
    type VARCHAR(20) DEFAULT 'regular' CHECK (type IN ('regular', 'vip', 'disabled', 'near_exit')),
    is_available BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_bus_seat UNIQUE (bus_id, number)
);

CREATE INDEX idx_seats_bus ON seats(bus_id);
COMMENT ON TABLE seats IS 'Места в автобусах';

-- ====================================
-- Таблица: users (Пользователи системы)
-- ====================================
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('cashier', 'dispatcher', 'controller', 'admin')),
    station_id UUID REFERENCES stations(id) ON DELETE SET NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_station ON users(station_id);
CREATE INDEX idx_users_role ON users(role);
COMMENT ON TABLE users IS 'Пользователи системы (кассиры, диспетчеры и т.д.)';

-- ====================================
-- Таблица: sessions (Сессии пользователей)
-- ====================================
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sessions_user ON sessions(user_id);
CREATE INDEX idx_sessions_token ON sessions(token_hash);
CREATE INDEX idx_sessions_expires ON sessions(expires_at);
COMMENT ON TABLE sessions IS 'Активные сессии пользователей';

-- ====================================
-- Таблица: tickets (Билеты)
-- ====================================
CREATE TABLE tickets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    trip_id UUID NOT NULL REFERENCES trips(id) ON DELETE RESTRICT,
    seat_id UUID REFERENCES seats(id) ON DELETE SET NULL,
    passenger_name VARCHAR(100),
    passport VARCHAR(20),
    phone VARCHAR(15),
    email VARCHAR(100),
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'returned', 'used')),
    return_penalty DECIMAL(10,2) DEFAULT 0,
    sold_at TIMESTAMP NOT NULL DEFAULT NOW(),
    sold_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    station_id UUID NOT NULL REFERENCES stations(id) ON DELETE RESTRICT,
    qr_code VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tickets_trip_seat ON tickets(trip_id, seat_id);
CREATE INDEX idx_tickets_status ON tickets(status);
CREATE INDEX idx_tickets_qr ON tickets(qr_code);
CREATE INDEX idx_tickets_passenger ON tickets(passenger_name);
COMMENT ON TABLE tickets IS 'Билеты';

-- ====================================
-- Таблица: blocking_rules (Правила блокировки мест)
-- ====================================
CREATE TABLE blocking_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    station_id UUID NOT NULL REFERENCES stations(id) ON DELETE CASCADE,
    route_id UUID REFERENCES routes(id) ON DELETE CASCADE,
    seat_range VARCHAR(20) NOT NULL,
    reason VARCHAR(50) NOT NULL CHECK (reason IN ('station_lock', 'privileged', 'maintenance')),
    valid_from DATE NOT NULL,
    valid_to DATE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT valid_date_range CHECK (valid_to >= valid_from)
);

CREATE INDEX idx_blocking_rules_route_station ON blocking_rules(route_id, station_id);
CREATE INDEX idx_blocking_rules_dates ON blocking_rules(valid_from, valid_to);
COMMENT ON TABLE blocking_rules IS 'Правила блокировки мест';

-- ====================================
-- Таблица: document_templates (Шаблоны документов)
-- ====================================
CREATE TABLE document_templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    format VARCHAR(10) NOT NULL CHECK (format IN ('pdf', 'docx')),
    is_default BOOLEAN DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_templates_default ON document_templates(is_default);
COMMENT ON TABLE document_templates IS 'Шаблоны документов (ПД-2 и др.)';
COMMENT ON COLUMN document_templates.content IS 'HTML с плейсхолдерами: {{trip_id}}, {{passengers}}, etc.';

-- ====================================
-- Таблица: audit_logs (Журнал аудита)
-- ====================================
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    action VARCHAR(50) NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    old_value JSONB,
    new_value JSONB,
    ip_address VARCHAR(45),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_user ON audit_logs(user_id);
CREATE INDEX idx_audit_created ON audit_logs(created_at);
COMMENT ON TABLE audit_logs IS 'Журнал аудита (152-ФЗ)';

-- ====================================
-- Таблица: boarding_events (События посадки)
-- ====================================
CREATE TABLE boarding_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    trip_id UUID NOT NULL REFERENCES trips(id) ON DELETE CASCADE,
    ticket_id UUID REFERENCES tickets(id) ON DELETE SET NULL,
    event_type VARCHAR(20) NOT NULL CHECK (event_type IN ('boarding_started', 'passenger_boarded')),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_boarding_trip ON boarding_events(trip_id);
CREATE INDEX idx_boarding_ticket ON boarding_events(ticket_id);
CREATE INDEX idx_boarding_type ON boarding_events(event_type);
COMMENT ON TABLE boarding_events IS 'События посадки пассажиров';

-- ====================================
-- Таблица: fiscal_receipts (Фискальные чеки)
-- ====================================
CREATE TABLE fiscal_receipts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ticket_id UUID REFERENCES tickets(id) ON DELETE SET NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('sale', 'refund')),
    amount DECIMAL(10,2) NOT NULL,
    ofd_url VARCHAR(255),
    kkt_serial VARCHAR(20),
    fiscal_document_number INTEGER,
    fiscal_sign VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_fiscal_ticket ON fiscal_receipts(ticket_id);
CREATE INDEX idx_fiscal_type ON fiscal_receipts(type);
CREATE INDEX idx_fiscal_created ON fiscal_receipts(created_at);
COMMENT ON TABLE fiscal_receipts IS 'Фискальные чеки (54-ФЗ)';

-- ====================================
-- Таблица: notifications (Уведомления)
-- ====================================
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type VARCHAR(20) NOT NULL CHECK (type IN ('sms', 'email', 'telegram', 'tts')),
    recipient VARCHAR(100) NOT NULL,
    message TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'sent', 'failed')),
    error_message TEXT,
    sent_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notifications_status ON notifications(status);
CREATE INDEX idx_notifications_created ON notifications(created_at);
COMMENT ON TABLE notifications IS 'Уведомления (SMS, Email, Telegram)';

-- ====================================
-- Таблица: announcements (Голосовые оповещения)
-- ====================================
CREATE TABLE announcements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    text TEXT NOT NULL,
    language VARCHAR(5) DEFAULT 'ru',
    priority VARCHAR(10) DEFAULT 'normal' CHECK (priority IN ('low', 'normal', 'high')),
    station_id UUID REFERENCES stations(id) ON DELETE CASCADE,
    played_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_announcements_station ON announcements(station_id);
CREATE INDEX idx_announcements_created ON announcements(created_at);
COMMENT ON TABLE announcements IS 'Голосовые оповещения (TTS)';

-- ====================================
-- Функции и триггеры
-- ====================================

-- Функция обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Триггеры updated_at для всех таблиц
CREATE TRIGGER update_stations_updated_at BEFORE UPDATE ON stations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_buses_updated_at BEFORE UPDATE ON buses
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_drivers_updated_at BEFORE UPDATE ON drivers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_routes_updated_at BEFORE UPDATE ON routes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_schedules_updated_at BEFORE UPDATE ON schedules
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_trips_updated_at BEFORE UPDATE ON trips
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_seats_updated_at BEFORE UPDATE ON seats
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tickets_updated_at BEFORE UPDATE ON tickets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_blocking_rules_updated_at BEFORE UPDATE ON blocking_rules
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_document_templates_updated_at BEFORE UPDATE ON document_templates
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ====================================
-- Начальные данные
-- ====================================

-- Создать системного администратора
-- Пароль: admin123 (хэш bcrypt)
INSERT INTO users (id, username, password_hash, full_name, role, is_active) VALUES
('00000000-0000-0000-0000-000000000001', 'admin', '$2a$10$8K1p/a0dL3LzW6R3b6V7JuDMKYJ0hPXkQkp6p3LN9Y8f0X2KF7Z3e', 'Системный администратор', 'admin', true);

COMMENT ON TABLE users IS 'Пользователи системы. Пароль admin по умолчанию: admin123';
