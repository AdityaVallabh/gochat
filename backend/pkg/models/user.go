package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name       string
	CreatedAt  time.Time
	LastActive time.Time
}
