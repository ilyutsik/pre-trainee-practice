package booking

type stubRooms struct {
	rooms map[string]Room
	fail  error
}

func (s *stubRooms) SaveRoom(r Room) error {
	s.rooms[r.ID] = r
	return nil
}

func (s *stubRooms) GetRoom(id string) (Room, error) {
	if s.fail != nil {
		return Room{}, s.fail
	}
	r, ok := s.rooms[id]
	if !ok {
		return Room{}, ErrNotFound
	}
	return r, nil
}

type stubReservations struct {
	*memReservations
	saveErr error
}

func (c *stubReservations) Save(r Reservation) error {
	if c.saveErr != nil {
		return c.saveErr
	}
	return c.memReservations.Save(r)
}
