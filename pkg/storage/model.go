package storage

type Order struct {
	ID       int    `db:"id"`
	OrderNum string `db:"order_num"`
	Data     []byte `db:"order_data"`
}
