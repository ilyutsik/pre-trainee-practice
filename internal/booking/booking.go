package booking

import (
	"errors"
	"fmt"
	"strings"
	"uuid"
)

var (
	ErrNotFound  = errors.New("booking: not found")
	ErrConflict  = errors.New("booking: slot conflict")
	ErrInvalid   = errors.New("booking: invalid argument")
	ErrCancelled = errors.New("booking: already cancelled")
)

type Room struct {
	ID   string
	Name string
}

type Reservation struct {
	ID        string
	RoomID    string
	Start     int // inclusive
	End       int // exclusive
	Cancelled bool
}

type RoomStore interface {
	SaveRoom(r Room) error
	GetRoom(id string) (Room, error)
}

type ReservationStore interface {
	Save(Reservation) error
	Get(id string) (Reservation, error)
	ListByRoom(roomID string) ([]Reservation, error)
}

type Service struct {
	rooms        RoomStore
	reservations ReservationStore
}

func NewService(rooms RoomStore, reservations ReservationStore) *Service {
	return &Service{rooms, reservations}
}

func (s *Service) AddRoom(id, name string) error {
	id, name = strings.TrimSpace(id), strings.TrimSpace(name)
	if id == "" {
		return fmt.Errorf("add room: id is required: %w", ErrInvalid)
	}
	if name == "" {
		return fmt.Errorf("add room: name is required: %w", ErrInvalid)
	}
	_, err := s.rooms.GetRoom(id)
	if err == nil {
		return fmt.Errorf("add room: duplicate id %s: %w", id, ErrInvalid)
	}
	if !errors.Is(err, ErrNotFound) {
		return fmt.Errorf("add room: %w", err)
	}
	if err := s.rooms.SaveRoom(Room{ID: id, Name: name}); err != nil {
		return fmt.Errorf("add room: save room: %w", err)
	}
	return nil
}

func (s *Service) Book(roomID string, start, end int) (Reservation, error) {
	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		return Reservation{}, fmt.Errorf("book room: id is required: %w", ErrInvalid)
	}
	if !(0 <= start) || !(start < end) || !(end <= 24) {
		return Reservation{}, fmt.Errorf("book room: invalid time interval %d-%d %w", start, end, ErrInvalid)
	}
	_, err := s.rooms.GetRoom(roomID)
	if err != nil {
		return Reservation{}, fmt.Errorf("book room: get room %s: %w", roomID, err)
	}
	list, err := s.reservations.ListByRoom(roomID)
	if err != nil {
		return Reservation{}, fmt.Errorf("book room: %w", err)
	}
	for _, reservation := range list {
		if !reservation.Cancelled {
			if start < reservation.End && end > reservation.Start {
				return Reservation{}, fmt.Errorf("book room: slot %d-%d conflicts with existing reservation: %s %w", start, end, reservation.ID, ErrConflict)
			}
		}
	}
	newUuid := uuid.New().String()
	newReservation := Reservation{ID: newUuid, RoomID: roomID, Start: start, End: end, Cancelled: false}

	err = s.reservations.Save(newReservation)
	if err != nil {
		return Reservation{}, fmt.Errorf("book 123 room: %w", err)
	}

	return newReservation, nil
}

func (s *Service) Cancel(reservationID string) error {
	reservationID = strings.TrimSpace(reservationID)
	if reservationID == "" {
		return fmt.Errorf("cancel: reservation id required: %w", ErrInvalid)
	}
	reservation, err := s.reservations.Get(reservationID)
	if err != nil {
		return fmt.Errorf("cancel: get reservation %s: %w", reservationID, err)
	}
	if reservation.Cancelled {
		return fmt.Errorf("cancel: reservation %s all ready canceled: %w", reservationID, ErrCancelled)
	}
	reservation.Cancelled = true
	err = s.reservations.Save(reservation)
	if err != nil {
		return fmt.Errorf("cancel: save reservation %s: %w", reservationID, err)
	}
	return nil
}

func (s *Service) Get(reservationID string) (Reservation, error) {
	reservationID = strings.TrimSpace(reservationID)
	if reservationID == "" {
		return Reservation{}, fmt.Errorf("get: reservation id required: %w", ErrInvalid)
	}

	reservation, err := s.reservations.Get(reservationID)
	if err != nil {
		return Reservation{}, fmt.Errorf("get: reservation %s: %w", reservationID, err)
	}

	return reservation, nil
}
func (s *Service) ListRoom(roomID string) ([]Reservation, error) {
	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		return []Reservation{}, fmt.Errorf("listRoom: room id required: %w", ErrInvalid)
	}
	if _, err := s.rooms.GetRoom(roomID); err != nil {
		return []Reservation{}, fmt.Errorf("list room: get room %s: %w", roomID, err)
	}
	reservationList, err := s.reservations.ListByRoom(roomID)
	if err != nil {
		return []Reservation{}, fmt.Errorf("listRoom: list room %s: %w", roomID, err)
	}

	return reservationList, nil
}
