package order

import (
	storage2 "L0/pkg/storage"
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/nats-io/stan.go"
	"log"
	"net/http"
	"strconv"
)

type Application struct {
	st *storage2.Storage
}

func NewApplication(st *storage2.Storage) *Application {
	return &Application{st: st}
}

func (a *Application) ByIDHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	param := ps.ByName("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		log.Print(err)
		response400(w)
		return
	}
	if id <= 0 {
		response400(w)
		return
	}

	order, err := a.st.OrderByID(id)
	if err != nil {
		log.Println(err)
		orderNotFound(w)
		return
	}
	fmt.Fprint(w, string(order))
}

func (a *Application) Consume(m *stan.Msg) {
	order := &storage2.Order{}
	order.Data = m.Data
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
