-- Создаём схему (как в LCM)
CREATE SCHEMA IF NOT EXISTS ltb;

-- Таблица пользователей
CREATE TABLE ltb.users (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone VARCHAR(20),
    role VARCHAR(20) NOT NULL DEFAULT 'client',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Таблица юристов
CREATE TABLE ltb.lawyers (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES ltb.users(id) ON DELETE CASCADE,
    full_name VARCHAR(200) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Таблица записей
CREATE TABLE ltb.appointments (
    id SERIAL PRIMARY KEY,
    client_id INTEGER NOT NULL REFERENCES ltb.users(id) ON DELETE CASCADE,
    lawyer_id INTEGER NOT NULL REFERENCES ltb.lawyers(id) ON DELETE CASCADE,
    appointment_time TIMESTAMP NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'confirmed',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Индексы для быстрого поиска
CREATE INDEX idx_appointments_lawyer_time ON ltb.appointments(lawyer_id, appointment_time);
CREATE INDEX idx_appointments_client ON ltb.appointments(client_id);
CREATE INDEX idx_users_telegram_id ON ltb.users(telegram_id);