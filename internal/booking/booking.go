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

type Service struct {
	rooms            map[string]string
	reservations     map[string]Reservation
	reservationOrder []string
}

func NewService() *Service {
	return &Service{rooms: make(map[string]string), reservations: make(map[string]Reservation), reservationOrder: make([]string, 0, 15)}
}

func (s *Service) AddRoom(id, name string) error {
	id, name = strings.TrimSpace(id), strings.TrimSpace(name)
	if id == "" {
		return fmt.Errorf("add room: id is required: %w", ErrInvalid)
	}
	if name == "" {
		return fmt.Errorf("add room: name is required: %w", ErrInvalid)
	}
	if _, exist := s.rooms[id]; exist {
		return fmt.Errorf("duplicate id %s: %w", id, ErrInvalid)
	}
	s.rooms[id] = name
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
	if _, exist := s.rooms[roomID]; !exist {
		return Reservation{}, fmt.Errorf("book room: room %s not found %w", roomID, ErrNotFound)
	}
	for _, id := range s.reservationOrder {
		reservation := s.reservations[id]
		if reservation.RoomID == roomID && !reservation.Cancelled {
			if start < reservation.End && end > reservation.Start {
				return Reservation{}, fmt.Errorf("book room: slot %d-%d conflicts with existing reservation: %s %w", start, end, reservation.ID, ErrConflict)
			}
		}
	}
	newUuid := uuid.New().String()
	s.reservationOrder = append(s.reservationOrder, newUuid)
	s.reservations[newUuid] = Reservation{newUuid, roomID, start, end, false}

	return s.reservations[newUuid], nil
}
func (s *Service) Cancel(reservationID string) error {
	reservationID = strings.TrimSpace(reservationID)
	if reservationID == "" {
		return fmt.Errorf("cancel: reservation id required: %w", ErrInvalid)
	}
	if _, exist := s.reservations[reservationID]; !exist {
		return fmt.Errorf("cancel: reservation %s not found: %w", reservationID, ErrNotFound)
	}
	if s.reservations[reservationID].Cancelled {
		return fmt.Errorf("cancel: reservation %s all ready canceled %w", reservationID, ErrCancelled)
	}
	res := s.reservations[reservationID]
	res.Cancelled = true
	s.reservations[reservationID] = res
	return nil
}

func (s *Service) Get(reservationID string) (Reservation, error) {
	reservationID = strings.TrimSpace(reservationID)
	if reservationID == "" {
		return Reservation{}, fmt.Errorf("get: reservation id required: %w", ErrInvalid)
	}
	if _, exist := s.reservations[reservationID]; !exist {
		return Reservation{}, fmt.Errorf("get: reservation %s not found: %w", reservationID, ErrNotFound)
	}

	return s.reservations[reservationID], nil
}
func (s *Service) ListRoom(roomID string) ([]Reservation, error) {
	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		return []Reservation{}, fmt.Errorf("listRoom: room id required: %w", ErrInvalid)
	}
	if _, exist := s.rooms[roomID]; !exist {
		return []Reservation{}, fmt.Errorf("listRoom: room with id %s not found: %w", roomID, ErrNotFound)
	}
	out := make([]Reservation, 0, len(s.reservations))
	for _, id := range s.reservationOrder {
		reservation := s.reservations[id]
		if reservation.RoomID == roomID {
			out = append(out, reservation)
		}
	}
	return out, nil
}
