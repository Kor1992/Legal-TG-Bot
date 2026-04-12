# Legal TG Bot

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-4169E1?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-✔-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![Telegram Bot API](https://img.shields.io/badge/Telegram%20Bot%20API-v5-26A5E4?style=flat&logo=telegram)](https://core.telegram.org/bots/api)

Telegram-бот для записи на юридические консультации. Позволяет клиентам регистрироваться, просматривать список доступных юристов, записываться на консультацию и получать информацию о своих записях.

## 📌 О проекте

Проект разработан в рамках обучения Go-разработке и является частью портфолио для стажировки в LegalTech.

**Основные возможности:**
- Регистрация клиентов с сохранением контактных данных
- Просмотр списка активных юристов
- Запись на консультацию с выбором доступного времени
- Просмотр своих записей
- Управление состоянием диалога (FSM)

## 🛠 Технологический стек

| Компонент | Технология |
|-----------|------------|
| Язык | Go 1.22+ |
| База данных | PostgreSQL 15+ |
| Драйвер БД | pgx/v5 |
| Telegram API | go-telegram-bot-api/v5 |
| Миграции | golang-migrate |
| Контейнеризация | Docker, Docker Compose |
| Конфигурация | .env, godotenv |

## 🏗 Архитектура

Проект следует принципам **чистой архитектуры** (Clean Architecture):


cmd/bot/ # Точка входа, DI
internal/
├── domain/ # Бизнес-модели (User, Lawyer, Appointment)
├── repository/ # Интерфейсы репозиториев
│ └── postgres/ # Реализации для PostgreSQL
├── service/ # Бизнес-логика
└── bot/ # Логика Telegram-бота
├── bot.go # Инициализация и главный цикл
├── handlers.go # Обработчики сообщений и команд
├── callbacks.go # Обработка inline-клавиатур
├── keyboards.go # Генерация клавиатур
├── states.go # Управление состояниями диалога
└── slots.go # Логика работы со слотами времени


**Почему чистая архитектура?**
- Разделение ответственности между слоями
- Лёгкая замена реализации (например, переход на другую БД)
- Удобство тестирования (зависимости через интерфейсы)

## 📦 База данных

### Схема

```sql
CREATE SCHEMA ltb;

-- Пользователи (клиенты и юристы)
CREATE TABLE ltb.users (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone VARCHAR(20),
    role VARCHAR(20) NOT NULL DEFAULT 'client',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Юристы (расширенная информация)
CREATE TABLE ltb.lawyers (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES ltb.users(id) ON DELETE CASCADE,
    full_name VARCHAR(200) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Записи на консультацию
CREATE TABLE ltb.appointments (
    id SERIAL PRIMARY KEY,
    client_id INTEGER NOT NULL REFERENCES ltb.users(id) ON DELETE CASCADE,
    lawyer_id INTEGER NOT NULL REFERENCES ltb.lawyers(id) ON DELETE CASCADE,
    appointment_time TIMESTAMP NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'confirmed',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);