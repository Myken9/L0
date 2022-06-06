package storage

type Order struct {
	ID      int    `db:"id"`
	OrderID string `db:"order_id"`
	Data    []byte `db:"order_data"`
}
