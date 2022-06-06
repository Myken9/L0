package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
)

type Storage struct {
	db Querier
	m  sync.Map
}

func NewStorage(conn Querier) *Storage {
	return &Storage{db: conn}
}

func (s *Storage) UploadCache(ctx context.Context) error {
	rows, err := s.db.Query(ctx, `SELECT * FROM orders ORDER BY id`)
	if err != nil {
		return fmt.Errorf("error getting orders from db: %w", err)
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(&order.ID, &order.OrderNum, &order.Data)
		if err != nil {
			log.Fatal(err)
		}
		orders = append(orders, order)
		s.m.Store(order.OrderNum, order.Data)
	}

	return nil
}

func (s *Storage) AddOrder(ctx context.Context, order *Order) error {
	q := "INSERT INTO orders (order_num, order_data) VALUES ($1, $2) RETURNING id"
	if err := s.db.QueryRow(ctx, q, order.OrderNum, order.Data).Scan(&order.ID); err != nil {
		return fmt.Errorf("error adding order: %w", err)
	}

	s.m.Store(order.OrderNum, order.Data)

	return nil
}
func (s *Storage) Order(orderNum string) ([]byte, error) {
	val, ok := s.m.Load(orderNum)
	if ok == false {
		log.Println("order not found in cache")
		return nil, errors.New("order not found")
	}
	value := val.([]byte)
	return value, nil
}
