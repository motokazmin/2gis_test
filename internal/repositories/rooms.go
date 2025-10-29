package repositories

import (
	"example/internal/dto"
	"example/internal/interfaces"
	"sync"
)

type roomsMemoryRepository struct {
	mu    sync.Mutex
	rooms []dto.Room
}

func NewRoomsMemoryRepository(knownRooms []dto.Room) interfaces.RoomsRepository {
	return &roomsMemoryRepository{
		rooms: knownRooms,
	}
}

func (o *roomsMemoryRepository) GetRoomsByHotel(hotelID string) ([]dto.Room, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	var rooms []dto.Room
	for _, room := range o.rooms {
		if room.HotelID == hotelID {
			rooms = append(rooms, room)
		}
	}
	return rooms, nil
}
