Скопируйте конфигурационный файл и поправьте, если нужно (по умолчанию прописаны переменные для подключения к докеру)
```bash
cp .env.dist .env
```

Поднимите контейнеры с nats-streaming + pgsql:
```bash
docker-compose up -d
```

Запустите продьюссера сообщений:
```bash
go run cmd/producer/main.go
```

Создайте таблицу со следующей структурой в контейнере (инструмент для миграций не использовал, создавал вручную)
```SQL
CREATE TABLE orders
(
    id BIGSERIAL PRIMARY KEY,
    order_num VARCHAR(255) NOT NULL UNIQUE,
    order_data jsonb NOT NULL
);
```

Запустите консьюмера:
```bash
go run cmd/consumer/main.go
```

> Примечание: заказом считается любое сообщение, у которого есть в JSON-е поле order_uid. Если необходимо другое поведение, 
> добавьте новые поля в order/app.go:39