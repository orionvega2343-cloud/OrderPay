# OrderPay - полноценный backend сервис учета заказов и оплаты, с DDD - архитектурой и кросс-атомарностью операций

## Стек
- Go 1.26
- Gin - HTTP-роутер
- PostgreSQL + sqlx - БД и работа с ней
- golang-migrate - миграции БД
- Swagger (swaggo) - документация API
- Docker / docker-compose - контейнеризация

## Сущности
Order - заказ пользователя
OrderItem - позиция в заказе
Payment - оплата пользователя
Refund - возврат

## Эндпоинты
    POST   /api/orders                - создать заказ
    GET    /api/orders                - получить все заказы
    GET    /api/orders/:id            - получить заказ по id
    PATCH  /api/orders/:id            - обновить статус заказа
    DELETE /api/orders/:id            - удалить заказ
    POST   /api/orders/:id/payments   - создать платеж по заказу (:id - id заказа)

    GET    /api/payments              - получить все платежи
    GET    /api/payments/total-amount - получить итоговую сумму платежей
    GET    /api/payments/:id          - получить платеж по id
    PUT    /api/payments/:id          - обновить платеж
    DELETE /api/payments/:id          - удалить платеж
    POST   /api/payments/:id/refunds  - создать возврат по платежу (:id - id платежа)

    GET    /api/refunds               - получить все возвраты

    GET    /api/payments/total-amount - принимает query параметры такие как: user_id, status, from, to

## Архитектура
Каждый домен разбит на:
- domain/ - бизнес-модели и правила(order.TransitionStatus - валидирует допустимые переходы статуса)
- repository/ - работа с БД
- service/ - бизнес логика
- handler/ - HTTP слой

## Unit of Work
Операция, затрагивающая несколько таблиц, обеспечивает кросс-атомарность через транзакцию,
выполняя SQL-запросы в рамках одной транзакции. Транзакция передается через context.Context,
репозитории сами решают использовать транзакцию или работать без нее
## Кросс-доменная атомарность
Пример: PaymentService.Create создаёт Payment и обновляет статус
связанного Order в рамках одной транзакции — если один из шагов
падает, откатываются оба.

## Перед запуском
Перед запуском необходимо создать `.env` в корне проекта. Пример полного файла:

```env
# пароль от БД (единственная обязательная переменная,
# читается из env, остальные поля БД - из configs/config_dev.yml)
DB_PASS=ваш_пароль

# опционально: порт БД на хосте для docker-compose,
# если 5432 уже занят локальным postgres
DB_HOST_PORT=5433
```

Host, порт БД внутри контейнера/сети, имя пользователя, имя БД и ssl_mode задаются
не через .env, а в `configs/config_dev.yml` (локальный запуск) и
`configs/config_docker.yml` (docker-compose) - поменяйте их там, если нужно.

## Запуск
    docker-compose up --build  - поднимает проект в изолированном контейнере
## Локальный запуск
    go build -o bin/orderpay ./cmd/server

