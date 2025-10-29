package errors

// ErrUnknownRoom - комната не найдена
type ErrUnknownRoom struct{}

func (e ErrUnknownRoom) Error() string {
	return "Room not found"
}

// ErrRoomUnavailable - номер не доступен на эти даты
type ErrRoomUnavailable struct{}

func (e ErrRoomUnavailable) Error() string {
	return "Room is not available for the selected dates"
}

// ErrFromAfterTo - ошибка валидации дат
type ErrFromAfterTo struct{}

func (e ErrFromAfterTo) Error() string {
	return "Invalid date range: 'from' must be before 'to'"
}
