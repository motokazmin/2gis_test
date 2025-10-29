package interfaces

import (
	"example/internal/dto"
)

type OrdersRepository interface {
	Create(order dto.Order) error
	GetOrdersByHotelAndRoom(hotelID, roomID string) ([]dto.Order, error)
}

type RoomsRepository interface {
	GetRoomsByHotel(hotelID string) ([]dto.Room, error)
}
