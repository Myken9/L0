package main

import (
	"L0/pkg/model"
	"encoding/json"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/nats-io/stan.go"
	"log"
	"time"
)

func main() {
	sc, err := stan.Connect(
		"test-cluster",
		"order_producer",
		stan.NatsURL("nats://localhost:4222"),
	)
	if err != nil {
		log.Println(err)

		return
	}
	defer sc.Close()

	t := time.NewTicker(time.Nanosecond * 1)
	for range t.C {
		log.Println("send msg to topic")

		if err = sc.Publish("order", genOrder()); err != nil {
			log.Fatal(err)
		}
	}
}

func genOrder() []byte {
	order := model.Order{
		OrderUid:          gofakeit.UUID(),
		TrackNumber:       gofakeit.UUID(),
		Entry:             gofakeit.UUID(),
		SmId:              gofakeit.Number(1, 100),
		Locale:            gofakeit.Language(),
		InternalSignature: gofakeit.UUID(),
		CustomerId:        gofakeit.UUID(),
		DeliveryService:   gofakeit.UUID(),
		Shardkey:          gofakeit.UUID(),
		OofShard:          gofakeit.UUID(),
		DateCreated:       gofakeit.Date(),

		Items: []model.Item{
			{
				ChrtId:      gofakeit.IntRange(1, 10000),
				TrackNumber: gofakeit.UUID(),
				Price:       gofakeit.IntRange(1, 10000),
				Rid:         gofakeit.UUID(),
				Name:        gofakeit.BeerName(),
				Sale:        0,
				Size:        "",
				TotalPrice:  0,
				NmId:        0,
				Brand:       "",
				Status:      0,
			},
		},
		Delivery: model.Delivery{
			Name:    gofakeit.StreetName(),
			Phone:   gofakeit.Phone(),
			Zip:     gofakeit.Zip(),
			City:    gofakeit.City(),
			Address: gofakeit.BitcoinAddress(),
			Region:  "",
			Email:   gofakeit.Email(),
		},
		Payment: model.Payment{
			Transaction:  gofakeit.UUID(),
			RequestId:    "",
			Currency:     "",
			Provider:     "",
			Amount:       0,
			PaymentDt:    0,
			Bank:         "",
			DeliveryCost: 0,
			GoodsTotal:   0,
			CustomFee:    0,
		},
	}

	b, err := json.Marshal(order)
	if err != nil {
		log.Fatal(err)
	}

	return b
}
