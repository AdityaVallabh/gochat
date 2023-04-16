package server

import (
	"log"
	"net/http"

	"github.com/AdityaVallabh/gochat/pkg/user"
	"github.com/google/uuid"
)

func (s *Server) handleWs() http.HandlerFunc {
	type request struct {
		UserID uuid.UUID
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := request{uuid.MustParse(r.URL.Query().Get("id"))}
		log.Println(req.UserID)
		v, ok := s.users.Load(req.UserID)
		if !ok {
			log.Println("no user")
			s.respond(w, r, http.StatusBadRequest, errorResponse{"no user"})
			return
		}
		u, ok := v.(*user.User)
		if !ok {
			log.Println("no user 2")
			s.respond(w, r, http.StatusBadRequest, errorResponse{"no user 2"})
			return
		}
		err := u.NewConn(w, r)
		if err != nil {
			log.Println("no join")
			s.respond(w, r, http.StatusBadRequest, errorResponse{err.Error()})
			return
		}
		log.Println("ok")
	}
}
