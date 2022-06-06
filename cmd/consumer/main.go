package main

import (
	"L0/order"
	"L0/pkg/storage"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
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

	store := storage.NewStorage(conn)
	if err := store.UploadCache(context.Background()); err != nil {
		log.Println("error uploading cache from db", err)
	}
	app := order.NewApplication(store)

	st, err := stan.Connect(
		"test-cluster",
		"order-consumer", stan.NatsURL(os.Getenv("NATS_STREAMING_URL")),
	)

	if err != nil {
		log.Fatal(err)
	}
	defer st.Close()

	if _, err = st.Subscribe("orders", app.Consume, stan.StartWithLastReceived()); err != nil {
		return
	}

	fmt.Println("Server is listening...")
	if err := http.ListenAndServe(os.Getenv("HTTP_ADDR"), configureRouter(app)); err != nil {
		log.Fatal(err)
	}
}

func configureRouter(app *order.Application) *httprouter.Router {
	router := httprouter.New()
	router.GET("/orders/:orderNum", app.ByIDHandler)
	return router
}
