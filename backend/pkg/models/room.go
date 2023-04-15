package models

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name      string
	CreatedAt time.Time
	DeletedAt time.Time `json:"-"`
}
