package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"example/internal/dto"
	apperrors "example/internal/errors"
	"example/internal/interfaces"
)

type OrderController struct {
	ordersService interfaces.OrdersService
}

func NewOrderController(ordersService interfaces.OrdersService) *OrderController {
	return &OrderController{ordersService: ordersService}
}

// Create создает новый заказ
func (c *OrderController) Create(w http.ResponseWriter, r *http.Request) {
	orderRequest, err := parseRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := c.ordersService.Create(orderRequest)
	if err != nil {
		// Различаем разные типы ошибок с помощью errors.As()
		var errFromAfterTo apperrors.ErrFromAfterTo
		var errUnknown apperrors.ErrUnknownRoom
		var errUnavailable apperrors.ErrRoomUnavailable

		if errors.As(err, &errFromAfterTo) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if errors.As(err, &errUnknown) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else if errors.As(err, &errUnavailable) {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = orderCreatedResponse(w, order)
	if err != nil {
		log.Default().Print(err)
	}
}

// Get получает заказ по ID (заглушка)
func (c *OrderController) Get(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// Update обновляет заказ (заглушка)
func (c *OrderController) Update(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// Delete удаляет заказ (заглушка)
func (c *OrderController) Delete(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func parseRequest(r *http.Request) (dto.Order, error) {
	var newOrder dto.Order
	if err := json.NewDecoder(r.Body).Decode(&newOrder); err != nil {
		return dto.Order{}, err
	}
	return newOrder, nil
}

func orderCreatedResponse(w http.ResponseWriter, order dto.Order) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(order)
	if err != nil {
		return err
	}
	return nil
}
