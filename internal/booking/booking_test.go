package booking

import (
	"errors"
	"io"
	"testing"
)

func TestBookConflict(t *testing.T) {
	s := NewMemoryService()
	if err := s.AddRoom("origin", "Six Seven"); err != nil {
		t.Fatal(err)
	}

	if _, err := s.Book("origin", 7, 10); err != nil {
		t.Fatal(err)
	}

	_, err := s.Book("origin", 8, 12)
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("got=%v; want=%v", err, ErrConflict)
	}
}

func TestBookAfterCancel(t *testing.T) {
	s := NewMemoryService()
	if err := s.AddRoom("origin", "Six Seven"); err != nil {
		t.Fatal(err)
	}

	reservation, err := s.Book("origin", 7, 10)
	if err != nil {
		t.Fatal(err)
	}

	if err := s.Cancel(reservation.ID); err != nil {
		t.Fatal(err)
	}
	_, err = s.Book("origin", 7, 10)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRoomStoreFailsOnGet(t *testing.T) {
	sr := &stubRooms{rooms: make(map[string]Room), fail: io.ErrClosedPipe}
	s := NewService(sr, newMemReservations())

	_, err := s.Book("origin", 9, 11)
	if !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("got=%v; want=%v", err, io.ErrClosedPipe)
	}
}

func TestReservationSaveFails(t *testing.T) {
	rooms := newMemRooms()
	err := rooms.SaveRoom(Room{ID: "origin", Name: "Six Seven"})
	if err != nil {
		t.Fatal(err)
	}

	res := newMemReservations()
	err = res.Save(Reservation{ID: "r1", RoomID: "origin", Start: 7, End: 10})
	if err != nil {
		t.Fatal(err)
	}

	s := NewService(rooms, &stubReservations{res, io.ErrClosedPipe})

	_, err = s.Book("origin", 12, 15)
	if !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("got=%v; want=%v", err, io.ErrClosedPipe)
	}
}

func TestCancelMissing(t *testing.T) {
	s := NewService(newMemRooms(), newMemReservations())

	err := s.Cancel("non present id")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got=%v; want=%v", err, ErrNotFound)
	}
}

func TestNotOverlapBook(t *testing.T) {
	s := NewMemoryService()
	err := s.AddRoom("origin", "Six Seven")
	if err != nil {
		t.Fatal(err)
	}

	_, err = s.Book("origin", 10, 12)
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.Book("origin", 12, 13)
	if err != nil {
		t.Fatal(err)
	}
}

func TestOverlapBook(t *testing.T) {
	cases := map[string]struct {
		s1, e1, s2, e2 int
		wantConflict   bool
	}{
		"overlap":        {10, 12, 11, 13, true},
		"touch at end":   {10, 12, 12, 14, false},
		"touch at start": {12, 13, 9, 12, false},
		"inside":         {10, 15, 12, 14, true},
	}
	s := NewMemoryService()
	bookId := "origin"
	for name, c := range cases {
		bookId = bookId + "1"
		t.Run(name, func(t *testing.T) {
			err := s.AddRoom(bookId, "Six Seven")
			if err != nil {
				t.Fatal(err)
			}
			_, err = s.Book(bookId, c.s1, c.e1)
			if err != nil {
				t.Fatal(err)
			}
			_, err = s.Book(bookId, c.s2, c.e2)
			if c.wantConflict {
				if !errors.Is(err, ErrConflict) {
					t.Errorf("got=%v; want=%v", err, ErrConflict)
				}
			} else if err != nil {
				t.Errorf("got=%v; want nil", err)
			}
		})

	}
}
