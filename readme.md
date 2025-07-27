# Subscriptions API

REST‑сервис для учёта онлайн‑подписок, реализованный на Go + Gin + PostgreSQL.

Функционал

- CRUD‑операции над подписками  
- Подсчёт суммарной стоимости по фильтру  
- Swagger‑документация  
- Миграции с Goose  
- Конфиг через `.env`  
- Запуск через Docker-compose
- Встроенный линтер
- Middleware для логирования

---

## Быстрый старт

1. Поднять контейнеры и мигрировать БД

   ```bash
   make build
   make up
   make goose-up
   ```

2. Открыть Swagger UI
   [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

## Использование CURL

```bash
# 1. Создать
curl -i -X POST http://localhost:8080/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025",
    "end_date": null
  }'

# 2. Получить
curl -i http://localhost:8080/subscriptions/<ID>

# 3. Обновить
curl -i -X PUT http://localhost:8080/subscriptions/<ID> \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus Premium",
    "price": 550,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025",
    "end_date": "12-2025"
  }'

# 4. Подсчитать сумму
curl -i -X POST http://localhost:8080/subscriptions/list \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "service_name": "Yandex Plus Premium",
    "start_date": "07-2025",
    "end_date": "12-2025"
  }'

# 5. Удалить
curl -i -X DELETE http://localhost:8080/subscriptions/<ID>
```

---

## Полезные команды

# Docker & сервисы
make build           # собрать Docker‑образы

make up              # поднять Postgres + API

make down            # остановить и удалить контейнеры и тома

make logs            # смотреть логи API

# Миграции Goose
make goose-add       # создать новую миграцию: goose-add NAME=<имя>

make goose-up        # применить все миграции

make goose-down      # откатить последнюю миграцию

make goose-status    # показать статус миграций

# Code quality
make lint            # запустить golangci‑lint