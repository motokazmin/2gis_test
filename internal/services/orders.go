package services

import (
	"example/internal/dto"
	"example/internal/errors"
	"example/internal/interfaces"
	"time"
)

type ordersService struct {
	ordersRepository      interfaces.OrdersRepository
	roomsMemoryRepository interfaces.RoomsRepository
}

func NewOrdersService(
	ordersRepository interfaces.OrdersRepository,
	roomsMemoryRepository interfaces.RoomsRepository,
) interfaces.OrdersService {
	return &ordersService{
		ordersRepository:      ordersRepository,
		roomsMemoryRepository: roomsMemoryRepository,
	}
}

func (o *ordersService) Create(order dto.Order) (dto.Order, error) {
	if err := o.validateCreateOrder(order); err != nil {
		return dto.Order{}, err
	}

	ok, err := o.isKnownRoom(order.HotelID, order.RoomTypeID)
	if err != nil {
		return dto.Order{}, err
	} else if !ok {
		return dto.Order{}, errors.ErrUnknownRoom{}
	}

	ok, err = o.isRoomAvailable(order.HotelID, order.RoomTypeID, order.From, order.To)
	if err != nil {
		return dto.Order{}, err
	} else if !ok {
		return dto.Order{}, errors.ErrRoomUnavailable{}
	}

	if err := o.ordersRepository.Create(order); err != nil {
		return dto.Order{}, err
	}

	return order, nil
}

func (o *ordersService) isRoomAvailable(hotelID string, roomID string, from time.Time, to time.Time) (bool, error) {
	orders, err := o.ordersRepository.GetOrdersByHotelAndRoom(hotelID, roomID)

	if err != nil {
		return false, err
	}

	if len(orders) == 0 {
		return true, nil
	}

	for _, order := range orders {
		if order.From.Before(to) && order.To.After(from) {
			return false, errors.ErrRoomUnavailable{}
		}
	}
	return true, nil
}

func (o *ordersService) isKnownRoom(hotelID string, roomID string) (bool, error) {
	rooms, err := o.roomsMemoryRepository.GetRoomsByHotel(hotelID)
	if err != nil {
		return false, err
	}
	for _, room := range rooms {
		if room.RoomTypeID == roomID {
			return true, nil
		}
	}
	return false, nil
}

func (o *ordersService) validateCreateOrder(order dto.Order) error {
	if order.From.IsZero() || order.To.IsZero() {
		return errors.ErrFromAfterTo{}
	}

	if order.From.After(order.To) {
		return errors.ErrFromAfterTo{}
	}

	return nil
}
