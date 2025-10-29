package app

import (
	"net/http"

	"example/internal/api/controllers"
	"example/internal/dto"
	httphandler "example/internal/http"
	"example/internal/interfaces"
	"example/internal/repositories"
	"example/internal/services"
)

const listenAddress = ":8080"

// App — оркестратор приложения, управляет всеми ресурсами
type App struct {
	server           *http.Server
	ordersRepository interfaces.OrdersRepository
	roomsRepository  interfaces.RoomsRepository
	ordersService    interfaces.OrdersService
	orderController  *controllers.OrderController
}

// NewApp инициализирует все компоненты приложения
func NewApp() *App {
	// Инициализируем репозитории
	ordersRepository := repositories.NewOrdersMemoryRepository()
	roomsRepository := repositories.NewRoomsMemoryRepository([]dto.Room{
		{"reddison", "lux"},
		{"reddison", "premium"},
	})

	// Инициализируем сервис
	ordersService := services.NewOrdersService(ordersRepository, roomsRepository)

	// Инициализируем контроллер
	orderController := controllers.NewOrderController(ordersService)

	// Инициализируем HTTP handler со всеми маршрутами и middleware
	handler := httphandler.NewHandler(orderController)

	// Создаём HTTP сервер
	server := &http.Server{
		Addr:    listenAddress,
		Handler: handler,
	}

	return &App{
		server:           server,
		ordersRepository: ordersRepository,
		roomsRepository:  roomsRepository,
		ordersService:    ordersService,
		orderController:  orderController,
	}
}

// Server возвращает HTTP сервер
func (a *App) Server() *http.Server {
	return a.server
}

// Close закрывает все ресурсы приложения
func (a *App) Close() error {
	// Когда будут реальные БД соединения, закроем их здесь
	// Например:
	// if closer, ok := a.ordersRepository.(io.Closer); ok {
	//     if err := closer.Close(); err != nil {
	//         return err
	//     }
	// }

	return nil
}
