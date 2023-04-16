package server

import (
	"net/http"
	"time"

	"github.com/AdityaVallabh/gochat/pkg/models"
	"github.com/AdityaVallabh/gochat/pkg/room"
	"github.com/AdityaVallabh/gochat/pkg/user"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func (s *Server) handleRoomPost() http.HandlerFunc {
	type response struct {
		ID uuid.UUID
	}

	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		rm := models.Room{
			CreatedAt: now,
		}
		err := s.DB.Save(&rm).Error
		if err != nil {
			log.Error(err.Error())
			s.respond(w, r, http.StatusInternalServerError, errorResponse{err.Error()})
			return
		}

		s.rooms.Store(rm.ID, room.NewRoom(&rm))

		log.WithFields(log.Fields{
			"id": rm.ID,
		}).Info("Created room")
		s.respond(w, r, http.StatusOK, response{
			ID: rm.ID,
		})
	}
}

func (s *Server) handleRoomJoinPost() http.HandlerFunc {
	type request struct {
		RoomID uuid.UUID
		UserID uuid.UUID
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		err := s.decode(w, r, &req)
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, errorResponse{err.Error()})
			return
		}
		v, ok := s.rooms.Load(req.RoomID)
		if !ok || v == nil {
			s.respond(w, r, http.StatusBadRequest, errorResponse{"no room"})
			return
		}
		rm, ok := v.(*room.Room)
		if !ok && rm != nil {
			s.respond(w, r, http.StatusBadRequest, errorResponse{"shouldnt happen"})
			return
		}

		v, ok = s.users.Load(req.UserID)
		if !ok {
			s.respond(w, r, http.StatusBadRequest, errorResponse{"no user"})
			return
		}
		u, ok := v.(*user.User)
		if !ok {
			s.respond(w, r, http.StatusBadRequest, errorResponse{"shddd user"})
			return
		}

		err = u.Join(rm)
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, errorResponse{"no join"})
			return
		}
		s.respond(w, r, http.StatusOK, rm)
	}
}
