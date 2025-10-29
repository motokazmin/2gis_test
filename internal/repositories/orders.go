package repositories

import (
	"example/internal/dto"
	"example/internal/interfaces"
	"sync"
)

type orderMemoryRepository struct {
	orders []dto.Order
	mu     sync.Mutex
}

func NewOrdersMemoryRepository() interfaces.OrdersRepository {
	return &orderMemoryRepository{
		orders: make([]dto.Order, 0),
	}
}

func (o *orderMemoryRepository) Create(order dto.Order) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.orders = append(o.orders, order)
	return nil
}

func (o *orderMemoryRepository) GetOrdersByHotelAndRoom(hotelID, roomID string) ([]dto.Order, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	var orders []dto.Order
	for _, order := range o.orders {
		if order.HotelID == hotelID && order.RoomTypeID == roomID {
			orders = append(orders, order)
		}
	}

	return orders, nil
}
