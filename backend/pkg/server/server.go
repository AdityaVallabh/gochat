package server

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/AdityaVallabh/gochat/pkg/models"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gorm.io/gorm"
)

type Server struct {
	Router *mux.Router
	DB     *gorm.DB

	rooms sync.Map // sync.Map[uuid.UUID]Room
	users sync.Map // sync.Map[uuid.UUID]User
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cors.Default().ServeHTTP(w, r, s.Router.ServeHTTP)
}

func (s *Server) Setup() error {
	err := s.DB.AutoMigrate(models.User{}, models.Room{})
	if err != nil {
		return fmt.Errorf("unable auto-migrate: %w", err)
	}
	s.routes()
	return nil
}
