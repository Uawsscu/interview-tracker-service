package entities

import (
	"time"

	"github.com/google/uuid"
)

type CardProgressLogs struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CardID    string    `json:"card_id"`
	ActorID   uuid.UUID `json:"actor_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
