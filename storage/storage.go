package storage

import (
	"L0/pkg/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"log"
	"sync"
)

type Queryer interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Ping(context.Context) error
	Prepare(context.Context, string, string) (*pgconn.StatementDescription, error)
}

type Storage struct {
	db Queryer
	mx sync.Mutex
	m  map[int][]byte
}

func NewStorage(conn Queryer) *Storage {
	return &Storage{db: conn, m: make(map[int][]byte)}
}

func (s *Storage) UploadCache(ctx context.Context) error {
	rows, err := s.db.Query(ctx, `
SELECT id, order_data 
FROM orders 
ORDER BY id`)
	if err != nil {
		return fmt.Errorf("error getting orders from db: %w", err)
	}
	defer rows.Close()

	var (
		order model.Order
		data  []byte
	)

	for rows.Next() {
		order = model.Order{}
		if err = scanOrder(&order, rows); err != nil {
			return fmt.Errorf("error fetching order: %w", err)
		}

		data, err = json.Marshal(order)
		if err != nil {
			return fmt.Errorf("error marshalling order: %w", err)
		}

		s.mx.Lock()
		s.m[order.Id] = data
		s.mx.Unlock()
	}

	log.Println("cache refreshed, count", len(s.m))

	return nil
}

func (s *Storage) AddOrder(ctx context.Context, order *model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("error marshalling order: %w", err)
	}

	err = s.db.QueryRow(ctx, `INSERT INTO orders (order_data) VALUES ($1) RETURNING id`, data).
		Scan(&order.Id)

	if err != nil {
		return fmt.Errorf("error adding order: %w", err)
	}

	s.mx.Lock()
	s.m[order.Id] = data
	s.mx.Unlock()

	return nil
}
func (s *Storage) GetOrderByID(id int) (order *model.Order, err error) {
	val, ok := s.m[id]
	if ok == false {
		log.Println("order not found in cache")
		return nil, err
	}
	if err = json.Unmarshal(val, &order); err != nil {
		log.Println("error unmarshalling order from cache: %w", err)
	}

	return order, nil
}

func scanOrder(o *model.Order, row pgx.Row) (err error) {
	if err = row.Scan(&o.Id, &o); err != nil {
		return fmt.Errorf("error scanning order id: %w", err)
	}

	return nil
}
