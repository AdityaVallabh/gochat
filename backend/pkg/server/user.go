package server

import (
	"net/http"
	"time"

	"github.com/AdityaVallabh/gochat/pkg/server/models"
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
		user := models.User{
			Name:       req.Name,
			CreatedAt:  now,
			LastActive: now,
		}
		s.DB.Save(&user)
		if err != nil {
			log.Error(err.Error())
			s.respond(w, r, http.StatusInternalServerError, errorResponse{err.Error()})
		}

		log.WithFields(log.Fields{
			"id":   user.ID,
			"name": user.Name,
		}).Info("Created user")
		s.respond(w, r, http.StatusOK, response{
			ID:   user.ID,
			Name: user.Name,
		})
	}
}
