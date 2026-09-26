# Backend Test — CS2 Items & User Balance API

Тестовое задание на позицию backend-разработчика. Два эндпоинта:
получение предметов Skinport с минимальными ценами и списание баланса пользователя.

## Стек

- Go 1.26.3
- PostgreSQL (без ORM, драйвер `pgx/v5`)
- Docker / Docker Compose

## Запуск

### Через Docker Compose

```bash
cp .env.example .env
docker-compose up --build
```

Приложение поднимется на `http://localhost:8080`, Postgres — на `localhost:5432`.

### Локально

```bash
cp .env.example .env
docker-compose up postgres   # поднять только БД
go run ./cmd/server
```

## ENV переменные

| Переменная      | Обязательна | Дефолт                                                    | Описание                          |
|-----------------|-------------|------------------------------------------------------------|------------------------------------|
| `HTTP_PORT`     | нет         | `8080`                                                     | Порт HTTP-сервера                  |
| `POSTGRES_DSN`  | да          | `postgres://app:app@localhost:5432/app?sslmode=disable`   | Строка подключения к Postgres      |
| `POSTGRES_USER` | да (для compose) | `app`                                                  | Юзер Postgres (только для контейнера БД) |
| `POSTGRES_PASSWORD` | да (для compose) | `app`                                              | Пароль Postgres (только для контейнера БД) |
| `POSTGRES_DB`   | да (для compose) | `app`                                                  | Имя БД (только для контейнера БД)  |

Полный пример — в `.env.example`.

## Эндпоинты

### `GET /items`

Список предметов Skinport с минимальными ценами отдельно для tradable и non-tradable лотов. Кешируется на 5 минут.

**Параметры (query):**
- `app_id` — опционально, дефолт `730` (CS2)
- `currency` — опционально, дефолт `EUR`, поддерживаются коды из [доки Skinport](https://docs.skinport.com/items)

**Пример:**
```bash
curl "http://localhost:8080/items?app_id=730&currency=EUR"
```

```json
[
  {
    "market_hash_name": "AK-47 | Redline (Field-Tested)",
    "currency": "EUR",
    "min_price_tradable": 10.5,
    "min_price_non_tradable": 8.0
  }
]
```

### `POST /users/{id}/withdraw`

Списывает сумму с баланса пользователя, атомарно записывает историю списания (было/стало/когда). Баланс не может уйти ниже нуля.

**Пример:**
```bash
curl -X POST http://localhost:8080/users/1/withdraw \
  -H "Content-Type: application/json" \
  -d '{"amount": "100.00"}'
```

**Ответ (200):**
```json
{
  "user_id": 1,
  "old_balance": "500",
  "new_balance": "400",
  "amount": "100",
  "created_at": "2026-09-26T12:00:00Z"
}
```

**Коды ошибок:**
- `400` — невалидные данные (некорректный `amount`, отсутствующий `id`)
- `404` — пользователь не найден
- `422` — недостаточно средств на балансе

В БД по умолчанию есть один пользователь: `id=1`, `balance=500.00`.

## Миграции

> **Важно:** сейчас миграция (`internal/migrations/0001_init.up.sql`) применяется через
> механизм `docker-entrypoint-initdb.d` официального Postgres-образа — он автоматически
> выполняет любые `*.sql`-файлы из примонтированной директории **только при первом
> старте** контейнера (на пустом volume).
>
> Это решение осознанно упрощено под рамки тестового задания. В реальном проекте
> для миграций стоит использовать полноценный инструмент — например
> [`goose`](https://github.com/pressly/goose) или [`golang-migrate`](https://github.com/golang-migrate/migrate) —
> чтобы иметь версионирование схемы, откаты (`down`-миграции) и возможность
> накатывать изменения на уже существующую БД без пересоздания контейнера.

## Тесты

```bash
go test ./...
```

## Структура проекта

```
cmd        — точка входа
internal/app       — сборка зависимостей, запуск/остановка сервера
internal/entity     — доменные модели и интерфейсы
internal/logic       — бизнес-логика
internal/repository   — реализации доступа к Postgres и Skinport API
internal/http         — HTTP-хендлеры и роутинг
internal/migrations   — SQL-миграции
pkg                    — переиспользуемые пакеты (логгер, HTTP-клиент, Postgres-пул, graceful shutdown)
```