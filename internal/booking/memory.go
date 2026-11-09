package booking

import (
	"fmt"
)

type memRooms struct {
	rooms map[string]Room
}

type memReservations struct {
	reservation      map[string]Reservation
	reservationOrder []string
}

func NewMemoryService() *Service {
	return NewService(newMemRooms(), newMemReservations())
}

func newMemRooms() *memRooms {
	return &memRooms{rooms: make(map[string]Room, 20)}
}

func (mem *memRooms) SaveRoom(r Room) error {
	mem.rooms[r.ID] = r
	return nil
}

func (mem *memRooms) GetRoom(id string) (Room, error) {
	room, exist := mem.rooms[id]
	if !exist {
		return Room{}, fmt.Errorf("%w: not found room with id=%v", ErrNotFound, id)
	}
	return room, nil
}

func newMemReservations() *memReservations {
	return &memReservations{reservation: make(map[string]Reservation, 20), reservationOrder: make([]string, 0, 20)}
}

func (mem *memReservations) Save(r Reservation) error {
	if _, exist := mem.reservation[r.ID]; !exist {
		mem.reservationOrder = append(mem.reservationOrder, r.ID)
	}
	mem.reservation[r.ID] = r
	mem.reservationOrder = append(mem.reservationOrder, r.ID)
	return nil
}

func (mem *memReservations) Get(id string) (Reservation, error) {
	reservation, exist := mem.reservation[id]
	if !exist {
		return Reservation{}, fmt.Errorf("%w: reservation with id %v is not exist", ErrNotFound, id)
	}
	return reservation, nil
}

func (mem *memReservations) ListByRoom(roomID string) ([]Reservation, error) {
	out := make([]Reservation, 0, 20)
	for _, id := range mem.reservationOrder {
		reservation := mem.reservation[id]
		if reservation.RoomID == roomID {
			out = append(out, reservation)
		}
	}
	return out, nil
}

var (
	_ RoomStore        = (*memRooms)(nil)
	_ ReservationStore = (*memReservations)(nil)
)
