package http

import (
	"net/http"

	"example/internal/api/controllers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewHandler инициализирует HTTP handler со всеми маршрутами и middleware
func NewHandler(orderController *controllers.OrderController) http.Handler {
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.Logger)

	// Маршруты для заказов
	router.Post("/orders", orderController.Create)

	return router
}
