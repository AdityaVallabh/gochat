package room

import (
	"errors"
	"sync"

	"github.com/AdityaVallabh/gochat/pkg/events"
	"github.com/AdityaVallabh/gochat/pkg/models"
	"github.com/AdityaVallabh/gochat/pkg/user"
	"github.com/google/uuid"
)

type Room struct {
	*models.Room

	mu    *sync.RWMutex
	ch    chan events.Message
	users map[uuid.UUID]*user.User
}

func (r *Room) ID() uuid.UUID {
	return r.Room.ID
}

func (r *Room) Add(u *user.User) error {
	defer r.mu.Unlock()
	r.mu.Lock()
	r.users[u.ID] = u
	return nil
}

func (r *Room) Remove(u *user.User) {
	defer r.mu.Unlock()
	r.mu.Lock()
	delete(r.users, u.ID)
}

func (r *Room) Speak(m events.Message) error {
	r.ch <- m
	return nil
}

func (r *Room) Delete() error {
	defer r.mu.Unlock()
	r.mu.Lock()
	if len(r.users) > 0 {
		return errors.New("cannot delete, users still using the room")
	}
	close(r.ch)
	return nil
}

func (r *Room) broadcastLoop() {
	for m := range r.ch {
		r.broadcast(m)
	}
}

func (r *Room) broadcast(m events.Message) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		u.Receive(m)
	}
}

func NewRoom(r *models.Room) *Room {
	room := &Room{
		Room:  r,
		users: make(map[uuid.UUID]*user.User),
		mu:    &sync.RWMutex{},
		ch:    make(chan events.Message),
	}
	go room.broadcastLoop()
	return room
}
