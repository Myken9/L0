package main

import (
	"L0/application"
	"L0/storage"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"github.com/nats-io/stan.go"
	"log"
	"net/http"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	if err = conn.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	stor := storage.NewStorage(conn)
	if err := stor.UploadCache(context.Background()); err != nil {
		log.Println("error uploading cache from db", err)
	}
	app := application.NewApplication(stor)

	st, err := stan.Connect(
		"test-cluster",
		"1", stan.NatsURL("nats://localhost:4222"),
	)

	if err != nil {
		log.Fatal(err)
	}
	defer st.Close()

	_, err = st.Subscribe("order", app.SubscribeToOrders, stan.DurableName("orders"))
	if err != nil {
		log.Println(err)

		return
	}

	msgHandler := msg("Hello")
	http.HandleFunc("/", app.InputOrderIDHandler)
	fmt.Println("Server is listening...")
	http.ListenAndServe("localhost:8181", msgHandler)
}

type msg string

func (m msg) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprint(resp, m)
}
