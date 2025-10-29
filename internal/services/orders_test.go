package services

import (
	"errors"
	"testing"
	"time"

	"example/internal/dto"
	apperrors "example/internal/errors"
	"example/internal/repositories"
)

func TestCreate_ValidOrder_Success(t *testing.T) {
	// Arrange
	ordersRepo := repositories.NewOrdersMemoryRepository()
	roomsRepo := repositories.NewRoomsMemoryRepository([]dto.Room{
		{"hotel1", "room1"},
	})
	service := NewOrdersService(ordersRepo, roomsRepo)

	now := time.Now()
	order := dto.Order{
		HotelID:    "hotel1",
		RoomTypeID: "room1",
		UserEmail:  "test@example.com",
		From:       now.AddDate(0, 0, 1),
		To:         now.AddDate(0, 0, 5),
	}

	// Act
	result, err := service.Create(order)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.HotelID != order.HotelID {
		t.Errorf("expected HotelID %s, got %s", order.HotelID, result.HotelID)
	}
}

func TestCreate_ZeroDates_ReturnError(t *testing.T) {
	// Arrange
	ordersRepo := repositories.NewOrdersMemoryRepository()
	roomsRepo := repositories.NewRoomsMemoryRepository([]dto.Room{})
	service := NewOrdersService(ordersRepo, roomsRepo)

	order := dto.Order{
		HotelID:    "hotel1",
		RoomTypeID: "room1",
		UserEmail:  "test@example.com",
		From:       time.Time{}, // Zero
		To:         time.Now(),
	}

	// Act
	_, err := service.Create(order)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var errFromAfterTo apperrors.ErrFromAfterTo
	if !errors.As(err, &errFromAfterTo) {
		t.Errorf("expected ErrFromAfterTo, got %T", err)
	}
}

func TestCreate_InvalidDateRange_ReturnError(t *testing.T) {
	// Arrange
	ordersRepo := repositories.NewOrdersMemoryRepository()
	roomsRepo := repositories.NewRoomsMemoryRepository([]dto.Room{})
	service := NewOrdersService(ordersRepo, roomsRepo)

	now := time.Now()
	order := dto.Order{
		HotelID:    "hotel1",
		RoomTypeID: "room1",
		UserEmail:  "test@example.com",
		From:       now.AddDate(0, 0, 10),
		To:         now.AddDate(0, 0, 5), // From > To
	}

	// Act
	_, err := service.Create(order)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var errFromAfterTo apperrors.ErrFromAfterTo
	if !errors.As(err, &errFromAfterTo) {
		t.Errorf("expected ErrFromAfterTo, got %T", err)
	}
}

func TestCreate_UnknownRoom_ReturnError(t *testing.T) {
	// Arrange
	ordersRepo := repositories.NewOrdersMemoryRepository()
	roomsRepo := repositories.NewRoomsMemoryRepository([]dto.Room{
		{HotelID: "hotel1", RoomTypeID: "room1"},
	})
	service := NewOrdersService(ordersRepo, roomsRepo)

	now := time.Now()
	order := dto.Order{
		HotelID:    "hotel1",
		RoomTypeID: "unknown_room", // Doesn't exist
		UserEmail:  "test@example.com",
		From:       now.AddDate(0, 0, 1),
		To:         now.AddDate(0, 0, 5),
	}

	// Act
	_, err := service.Create(order)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var errUnknown apperrors.ErrUnknownRoom
	if !errors.As(err, &errUnknown) {
		t.Errorf("expected ErrUnknownRoom, got %T", err)
	}
}

func TestCreate_RoomUnavailable_ReturnError(t *testing.T) {
	// Arrange
	ordersRepo := repositories.NewOrdersMemoryRepository()
	roomsRepo := repositories.NewRoomsMemoryRepository([]dto.Room{
		{"hotel1", "room1"},
	})
	service := NewOrdersService(ordersRepo, roomsRepo)

	now := time.Now()

	// Create first order
	firstOrder := dto.Order{
		HotelID:    "hotel1",
		RoomTypeID: "room1",
		UserEmail:  "test1@example.com",
		From:       now.AddDate(0, 0, 1),
		To:         now.AddDate(0, 0, 5),
	}
	service.Create(firstOrder)

	// Try to create overlapping order
	secondOrder := dto.Order{
		HotelID:    "hotel1",
		RoomTypeID: "room1",
		UserEmail:  "test2@example.com",
		From:       now.AddDate(0, 0, 3), // Overlaps with first
		To:         now.AddDate(0, 0, 7),
	}

	// Act
	_, err := service.Create(secondOrder)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var errUnavailable apperrors.ErrRoomUnavailable
	if !errors.As(err, &errUnavailable) {
		t.Errorf("expected ErrRoomUnavailable, got %T", err)
	}
}

func TestRoomAvailable_NoOrders_ReturnTrue(t *testing.T) {
	// Arrange
	ordersRepo := repositories.NewOrdersMemoryRepository()
	roomsRepo := repositories.NewRoomsMemoryRepository([]dto.Room{})
	service := &ordersService{
		ordersRepository: ordersRepo,
		roomsRepository:  roomsRepo,
	}

	now := time.Now()

	// Act
	available, err := service.isRoomAvailable("hotel1", "room1", now.AddDate(0, 0, 1), now.AddDate(0, 0, 5))

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !available {
		t.Error("expected room to be available")
	}
}

func TestRoomAvailable_NonIntersectingOrders_ReturnTrue(t *testing.T) {
	// Arrange
	ordersRepo := repositories.NewOrdersMemoryRepository()
	roomsRepo := repositories.NewRoomsMemoryRepository([]dto.Room{})
	service := &ordersService{
		ordersRepository: ordersRepo,
		roomsRepository:  roomsRepo,
	}

	now := time.Now()

	// Create order: 1-5 Jan
	existingOrder := dto.Order{
		HotelID:    "hotel1",
		RoomTypeID: "room1",
		From:       now.AddDate(0, 0, 1),
		To:         now.AddDate(0, 0, 5),
	}
	ordersRepo.Create(existingOrder)

	// Check availability: 6-10 Jan (no overlap)
	available, err := service.isRoomAvailable("hotel1", "room1", now.AddDate(0, 0, 6), now.AddDate(0, 0, 10))

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !available {
		t.Error("expected room to be available for non-overlapping dates")
	}
}

func TestRoomAvailable_IntersectingOrders_ReturnFalse(t *testing.T) {
	// Arrange
	ordersRepo := repositories.NewOrdersMemoryRepository()
	roomsRepo := repositories.NewRoomsMemoryRepository([]dto.Room{})
	service := &ordersService{
		ordersRepository: ordersRepo,
		roomsRepository:  roomsRepo,
	}

	now := time.Now()

	// Create order: 1-5 Jan
	existingOrder := dto.Order{
		HotelID:    "hotel1",
		RoomTypeID: "room1",
		From:       now.AddDate(0, 0, 1),
		To:         now.AddDate(0, 0, 5),
	}
	ordersRepo.Create(existingOrder)

	// Check availability: 3-7 Jan (overlaps!)
	available, err := service.isRoomAvailable("hotel1", "room1", now.AddDate(0, 0, 3), now.AddDate(0, 0, 7))

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if available {
		t.Error("expected room to be unavailable for overlapping dates")
	}
}

func TestKnownRoom_RoomExists_ReturnTrue(t *testing.T) {
	// Arrange
	ordersRepo := repositories.NewOrdersMemoryRepository()
	roomsRepo := repositories.NewRoomsMemoryRepository([]dto.Room{
		{"hotel1", "room1"},
	})
	service := &ordersService{
		ordersRepository: ordersRepo,
		roomsRepository:  roomsRepo,
	}

	// Act
	known, err := service.isKnownRoom("hotel1", "room1")

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !known {
		t.Error("expected room to be known")
	}
}

func TestKnownRoom_RoomNotExists_ReturnFalse(t *testing.T) {
	// Arrange
	ordersRepo := repositories.NewOrdersMemoryRepository()
	roomsRepo := repositories.NewRoomsMemoryRepository([]dto.Room{
		{"hotel1", "room1"},
	})
	service := &ordersService{
		ordersRepository: ordersRepo,
		roomsRepository:  roomsRepo,
	}

	// Act
	known, err := service.isKnownRoom("hotel1", "unknown_room")

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if known {
		t.Error("expected room to be unknown")
	}
}

func TestValidateCreateOrder_ValidOrder_NoError(t *testing.T) {
	// Arrange
	ordersRepo := repositories.NewOrdersMemoryRepository()
	roomsRepo := repositories.NewRoomsMemoryRepository([]dto.Room{})
	service := &ordersService{
		ordersRepository: ordersRepo,
		roomsRepository:  roomsRepo,
	}

	now := time.Now()
	order := dto.Order{
		From: now.AddDate(0, 0, 1),
		To:   now.AddDate(0, 0, 5),
	}

	// Act
	err := service.validateCreateOrder(order)

	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestValidateCreateOrder_FromAfterTo_ReturnError(t *testing.T) {
	// Arrange
	ordersRepo := repositories.NewOrdersMemoryRepository()
	roomsRepo := repositories.NewRoomsMemoryRepository([]dto.Room{})
	service := &ordersService{
		ordersRepository: ordersRepo,
		roomsRepository:  roomsRepo,
	}

	now := time.Now()
	order := dto.Order{
		From: now.AddDate(0, 0, 5),
		To:   now.AddDate(0, 0, 1), // From > To
	}

	// Act
	err := service.validateCreateOrder(order)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var errFromAfterTo apperrors.ErrFromAfterTo
	if !errors.As(err, &errFromAfterTo) {
		t.Errorf("expected ErrFromAfterTo, got %T", err)
	}
}
