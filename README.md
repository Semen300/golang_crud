# crud-go

`crud-go` — это REST API на Go для магазинов/заказов с ролью клиента, рабочего и менеджера. Проект использует Gin, PostgreSQL и JWT-авторизацию.

## Основные возможности

- Регистрация и аутентификация пользователей
- JWT access / refresh токены
- Отдельные роли: customer, worker, manager
- Управление заказами, корзиной, задачами и товарами
- Настраиваемая логика работы с заказами и назначением исполнителя
- Тестирование с использованием `stretchr/testify` и `testcontainers`

## Архитектура

- `cmd/app/main.go` — точка входа приложения
- `internal/handler` — HTTP-обработчики
- `internal/service` — бизнес-логика
- `internal/repository` — подключение к БД и CRUD операции
- `internal/middleware` — промежуточная авторизация
- `internal/model` — модели данных

## Технологии

- Go 1.26
- Gin Web Framework
- PostgreSQL (pgx)
- JWT
- dotenv

## Требования

- Go 1.26+
- PostgreSQL
- Переменные окружения, описанные ниже

## Переменные окружения

Пример `.env`:

```env
APP_PORT=8080
DATABASE_USERNAME=postgres
DATABASE_PASSWORD=password
DATABASE_PORT=5432
DATABASE_SCHEMA=crud_db
AUTH_ACCESS_LIFETIME=300000 //Время жизни access-токена в миллисекундах
AUTH_REFRESH_LIFETIME=86400000 //Время жизни refresh-токена в миллисекундах
AUTH_ACCESS_KEY=your_access_secret
AUTH_REFRESH_KEY=your_refresh_secret
```

## Запуск

1. Установите PostgreSQL.
2. Создайте базу данных, указанную в `DATABASE_SCEMA`.
3. Создайте `.env` с нужными переменными.
4. Запустите приложение:

```bash
go run ./cmd/app
```

## API

### Auth

- `POST /api/v1/auth/login` — получить access/refresh токены
- `POST /api/v1/auth/refresh` — обновить access токен
- `POST /api/v1/auth/logout` — выйти

### Customer

- `POST /api/v1/customer` — регистрация клиента
- `GET /api/v1/customer/orders` — список заказов
- `GET /api/v1/customer/orders/:id` — заказ по ID
- `POST /api/v1/customer/orders` — создать заказ
- `DELETE /api/v1/customer/orders/:id` — удалить заказ
- `GET /api/v1/customer/items` — список товаров
- `GET /api/v1/customer/basket` — корзина
- `POST /api/v1/customer/basket` — добавить товар в корзину
- `DELETE /api/v1/customer/basket/:id` — удалить товар из корзины

### Worker

- `GET /api/v1/worker/orders` — заказы рабочего
- `GET /api/v1/worker/orders/:id` — заказ рабочего по ID
- `POST /api/v1/worker/tasks/:id` — отметить задачу как выполненную

### Manager

- `GET /api/v1/manager/orders` — заказы менеджера
- `GET /api/v1/manager/orders/:id` — заказ по ID
- `POST /api/v1/manager/orders/:id` — назначить рабочего на заказ
- `GET /api/v1/manager/workers` — список рабочих

## Тесты

```bash
go test ./...
```

## Покрытие тестами

Реализовано покрытие 75% критического функционала Unit- и интеграционными тестами. Файлы покрытия тестами можно посмотреть в директории coverage.

## Примечания

- Сервис авторизации реализует хранение и отзыв refresh токенов.
- Все защищенные маршруты используют middleware `AuthMiddlewareFunc`.
