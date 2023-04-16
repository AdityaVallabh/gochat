package user

import (
	"errors"
	"net/http"
	"sync"

	"github.com/AdityaVallabh/gochat/pkg/events"
	"github.com/AdityaVallabh/gochat/pkg/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var (
	clientDisconnects = []int{websocket.CloseGoingAway, websocket.CloseAbnormalClosure}

	ErrUserNotInRoom = errors.New("user not in room")
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type room interface {
	ID() uuid.UUID
	Add(*User) error
	Remove(*User) error
	Send(events.Message) error
}

type User struct {
	*models.User
	Conns    map[uuid.UUID]*websocket.Conn
	MaxConns int
	Rooms    sync.Map // sync.Map[uuid.UUID]room

	mu  *sync.RWMutex
	ch  chan events.Message
	log *log.Entry
}

func (u *User) NewConn(w http.ResponseWriter, r *http.Request) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if len(u.Conns) == u.MaxConns {
		return errors.New("max conns reached")
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		w.Write([]byte(err.Error()))
		return nil
	}
	id := uuid.New()
	u.Conns[id] = ws
	go u.reader(id, ws)
	return nil
}

func (u *User) Join(r room) error {
	err := r.Add(u)
	if err != nil {
		return err
	}
	u.Rooms.Store(r.ID(), r)
	return nil
}

func (u *User) Leave(id uuid.UUID) error {
	v, ok := u.Rooms.LoadAndDelete(id)
	if !ok {
		return ErrUserNotInRoom
	}
	r, ok := v.(room)
	if !ok {
		return errors.New("unable to load room")
	}
	return r.Remove(u)
}

func (u *User) Receive(m events.Message) error {
	u.ch <- m
	return nil
}

func (u *User) sendLoop() {
	for m := range u.ch {
		u.send(m)
	}
}

func (u *User) send(m events.Message) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	for _, conn := range u.Conns {
		go func(conn *websocket.Conn) {
			err := conn.WriteJSON(m)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err.Error(),
				}).Error("could not send message")
			}
		}(conn)
	}
}

func (u *User) reader(id uuid.UUID, ws *websocket.Conn) {
	defer func() {
		u.mu.Lock()
		defer u.mu.Unlock()
		delete(u.Conns, id)
		log.Info("Client disconnected")
	}()
	log.Info("Client connected")
	for {
		var m events.Message
		err := ws.ReadJSON(&m)
		if err != nil {
			log := log.WithField("error", err)
			if websocket.IsUnexpectedCloseError(err, clientDisconnects...) {
				log.Warn("error reading message from ws")
				return
			}
			if websocket.IsCloseError(err, clientDisconnects...) {
				return
			}
			log.Warn("error decoding message")
			continue
		}
		err = u.speak(m)
		if err != nil {
			log.Error(err)
		}
	}
}

func (u *User) speak(m events.Message) error {
	v, ok := u.Rooms.Load(m.Room.ID)
	if !ok {
		return ErrUserNotInRoom
	}
	r, ok := v.(room)
	if !ok {
		return errors.New("unable to load room")
	}
	return r.Send(m)
}

func NewUser(u *models.User) *User {
	user := &User{
		User:     u,
		MaxConns: 2,
		Rooms:    sync.Map{},
		Conns:    make(map[uuid.UUID]*websocket.Conn),
		mu:       &sync.RWMutex{},
		ch:       make(chan events.Message),
		log:      log.WithField("user", u.Name),
	}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	go user.sendLoop()
	return user
}
