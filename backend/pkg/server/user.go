package server

import (
	"net/http"
	"time"

	"github.com/AdityaVallabh/gochat/pkg/models"
	"github.com/AdityaVallabh/gochat/pkg/user"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func (s *Server) handleUserPost() http.HandlerFunc {
	type request struct {
		Name string
	}
	type response struct {
		ID   uuid.UUID
		Name string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		err := s.decode(w, r, &req)
		if err != nil {
			log.Warn(err.Error())
			s.respond(w, r, http.StatusBadRequest, errorResponse{err.Error()})
			return
		}

		now := time.Now()
		u := models.User{
			Name:       req.Name,
			CreatedAt:  now,
			LastActive: now,
		}
		s.DB.Save(&u)
		if err != nil {
			log.Error(err.Error())
			s.respond(w, r, http.StatusInternalServerError, errorResponse{err.Error()})
		}
		s.users.Store(u.ID, user.NewUser(&u))
		log.WithFields(log.Fields{
			"id":   u.ID,
			"name": u.Name,
		}).Info("Created user")
		s.respond(w, r, http.StatusOK, response{
			ID:   u.ID,
			Name: u.Name,
		})
	}
}
