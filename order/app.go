package order

import (
	"L0/pkg/storage"
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/nats-io/stan.go"
	"log"
	"net/http"
)

type Application struct {
	st *storage.Storage
}

func NewApplication(st *storage.Storage) *Application {
	return &Application{st: st}
}

func (a *Application) ByIDHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	orderNum := ps.ByName("orderNum")
	if orderNum == "" {
		response400(w)
		return
	}

	order, err := a.st.Order(orderNum)
	if err != nil {
		log.Println(err)
		orderNotFound(w)
		return
	}
	fmt.Fprint(w, string(order))
}

type natsMessage struct {
	OrderNum string `json:"order_uid"`
}

func (a *Application) Consume(m *stan.Msg) {
	var msg natsMessage

	// Check if message has field "order_uid"
	if err := json.Unmarshal(m.Data, &msg); err != nil {
		log.Println(err)
		return
	}

	order := &storage.Order{}
	order.Data = m.Data
	order.OrderNum = msg.OrderNum

	if err := a.st.AddOrder(context.Background(), order); err != nil {
		log.Println(err)
	}
}

func response400(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, http.StatusText(http.StatusBadGateway))
}

func orderNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "order not found")
}
