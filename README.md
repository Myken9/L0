

```bash
cp .env.dist .env
```

Создайте таблицу со следующей структурой
```SQL
CREATE TABLE orders
(
    id BIGSERIAL PRIMARY KEY,
    order_data jsonb NOT NULL
);
```