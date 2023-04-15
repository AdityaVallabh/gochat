package user

import (
	"errors"
	"testing"

	"github.com/AdityaVallabh/gochat/pkg/events"
	"github.com/AdityaVallabh/gochat/pkg/models"
	"github.com/google/uuid"
)

func TestUser_JoinSpeakLeave(t *testing.T) {
	u := NewUser(&models.User{
		Name: "user1",
	})
	r := newTestRoom()
	t.Run("user should be able to join a room", func(t *testing.T) {
		if err := u.Join(&r); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if err := u.speak(events.Message{Room: &models.Room{ID: r.id}}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if err := u.speak(events.Message{Room: &models.Room{ID: r.id}}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if r.msgs != 2 {
			t.Errorf("expected r1.msgs to be %v, got %v", 2, r.msgs)
		}
		if err := u.Leave(r.id); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if err := u.speak(events.Message{Room: &models.Room{ID: r.id}}); !errors.Is(err, ErrUserNotInRoom) {
			t.Errorf("expected error after leave+speak but got %v", err)
		}
		if r.msgs != 2 {
			t.Errorf("expected r1.msgs to be %v, got %v", 2, r.msgs)
		}
	})
}

type testRoom struct {
	id   uuid.UUID
	msgs int
}

func (r testRoom) ID() uuid.UUID {
	return r.id
}

func (r testRoom) Add(_ *User) error {
	return nil
}

func (r testRoom) Remove(_ *User) error {
	return nil
}

func (r *testRoom) Send(_ events.Message) error {
	r.msgs++
	return nil
}

func newTestRoom() testRoom {
	return testRoom{
		id: uuid.New(),
	}
}
