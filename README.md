

```bash
cp .env.dist .env
```

Создайте таблицу со следующей структурой
```SQL
CREATE TABLE orders
(
    id BIGSERIAL PRIMARY KEY,
    order_num VARCHAR(255) NOT NULL UNIQUE,
    order_data jsonb NOT NULL
);
```