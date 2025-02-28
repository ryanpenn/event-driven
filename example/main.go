package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"event-driven-system/cqrs/command"
	"event-driven-system/cqrs/model"
	"event-driven-system/cqrs/projection"
	"event-driven-system/cqrs/query"
)

type HttpHandler func(http.ResponseWriter, *http.Request) error

func MakeHandler(fn HttpHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("request", r.Method, r.URL.Path)
		if err := fn(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func orderCreate(w http.ResponseWriter, r *http.Request) error {
	var order model.Order
	cmd := command.NewCreateOrderCommand(&order)
	if err := cmd.Execute(r.Context()); err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(fmt.Appendf(nil, "Order created:%s", order.ID))
	return nil
}

func orderUpdate(w http.ResponseWriter, r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	var order model.Order
	if err := json.Unmarshal(body, &order); err != nil {
		return err
	}
	cmd := command.NewUpdateOrderCommand(&order)
	if err := cmd.Execute(r.Context()); err != nil {
		return err
	}
	w.WriteHeader(http.StatusAccepted)
	return nil
}

func orderDelete(w http.ResponseWriter, r *http.Request) error {
	orderId := r.URL.Query().Get("id")
	cmd := command.NewDeleteOrderCommand(orderId)
	if err := cmd.Execute(r.Context()); err != nil {
		return err
	}
	w.WriteHeader(http.StatusAccepted)
	return nil
}

func orderGet(w http.ResponseWriter, r *http.Request) error {
	orderId := r.URL.Query().Get("id")
	order, err := query.OrderQueryService{}.GetOrder(r.Context(), orderId)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
	return nil
}

func orderList(w http.ResponseWriter, r *http.Request) error {
	orders, err := query.OrderQueryService{}.ListOrders(r.Context())
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
	return nil
}

func orderStatus(w http.ResponseWriter, r *http.Request) error {
	orderId := r.URL.Query().Get("id")
	orderStatus, err := query.OrderQueryService{}.GetOrderStatus(r.Context(), orderId)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(orderStatus))
	return nil
}

func listEvent(w http.ResponseWriter, r *http.Request) error {
	streamID := r.URL.Query().Get("id")
	w.Header().Set("Content-Type", "application/json")
	events, err := query.EventQueryService{}.ListEvents(r.Context(), streamID)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
	return nil
}

func initData() {
	cmd := command.NewCreateOrderCommand(&model.Order{
		ID:     "1",
		Status: "created",
		Items: []model.OrderItem{
			{ID: "1", Name: "item1", Quantity: 1, Price: 10.0},
			{ID: "2", Name: "item2", Quantity: 2, Price: 20.0},
		},
	})
	if err := cmd.Execute(context.TODO()); err != nil {
		fmt.Println(err)
	}
}

func main() {
	b := projection.NewOrderProjectionBuilder()
	b.Init()
	defer b.Close()

	initData()

	http.HandleFunc(fmt.Sprintf("%s /order/create", http.MethodPost), MakeHandler(orderCreate))
	http.HandleFunc(fmt.Sprintf("%s /order/update", http.MethodPut), MakeHandler(orderUpdate))
	http.HandleFunc(fmt.Sprintf("%s /order/delete", http.MethodDelete), MakeHandler(orderDelete))
	http.HandleFunc(fmt.Sprintf("%s /order/get", http.MethodGet), MakeHandler(orderGet))
	http.HandleFunc(fmt.Sprintf("%s /order/list", http.MethodGet), MakeHandler(orderList))
	http.HandleFunc(fmt.Sprintf("%s /order/status", http.MethodGet), MakeHandler(orderStatus))
	http.HandleFunc(fmt.Sprintf("%s /event/list", http.MethodGet), MakeHandler(listEvent))

	fmt.Println("Server started at :9000")
	http.ListenAndServe(":9000", nil)
}
