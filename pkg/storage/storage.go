package storage

import (
	"context"
	"encoding/json"
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
		err := rows.Scan(&order.ID, &order.Data)
		if err != nil {
			log.Fatal(err)
		}
		orders = append(orders, order)
		s.m.Store(order.ID, order.Data)
	}

	return nil
}

func (s *Storage) AddOrder(ctx context.Context, order *Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("error marshalling order: %w", err)
	}

	err = s.db.QueryRow(ctx, `INSERT INTO orders (order_data) VALUES ($1) RETURNING id`, data).
		Scan(&order.ID)

	if err != nil {
		return fmt.Errorf("error adding order: %w", err)
	}

	s.m.Store(order.ID, data)

	return nil
}
func (s *Storage) OrderByID(id int) ([]byte, error) {
	val, ok := s.m.Load(id)
	if ok == false {
		log.Println("order not found in cache")
		return nil, errors.New("order not found")
	}
	value := val.([]byte)
	return value, nil
}
