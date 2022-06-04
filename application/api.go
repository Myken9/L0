package application

import (
	"L0/pkg/model"
	"L0/storage"
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Application struct {
	st *storage.Storage
}

func NewApplication(st *storage.Storage) *Application {
	return &Application{st: st}
}

var data = "Go Template"
var tmpl, _ = template.New("data").Parse("<h1>{{ .}}</h1>")

func (a *Application) InputOrderIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("orderid")
	if id == "" {
		tmpl.Execute(w, nil)

		return
	}

	iid, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	order, err := a.st.GetOrderByID(iid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	tmpl.Execute(w, order)
}

func (a *Application) SubscribeToOrders(m *stan.Msg) {
	ord, err := orderFromJSON(m.Data)
	if err != nil {
		log.Println("error:", err)

		return
	}

	if err = a.st.AddOrder(context.Background(), ord); err != nil {
		log.Println(err)
	}
}

func orderFromJSON(data []byte) (*model.Order, error) {
	order := &model.Order{}

	if err := json.Unmarshal(data, order); err != nil {
		return nil, fmt.Errorf("error unmarshalling order: %w", err)
	}

	return order, nil
}
