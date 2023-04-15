package events

import (
	"time"

	"github.com/AdityaVallabh/gochat/pkg/models"
)

type Message struct {
	Room *models.Room
	User *models.User
	Time time.Time
	Data string
}
